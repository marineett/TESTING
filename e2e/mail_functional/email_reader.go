package service_logic

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// IEmailReader интерфейс для чтения email
type IEmailReader interface {
	GetLastMessage() (string, error) // Возвращает содержимое последнего сообщения
}

// EmailReader сервис для чтения email через IMAP
type EmailReader struct {
	imapHost     string
	imapPort     string
	imapUser     string
	imapPassword string
}

// CreateEmailReader создает сервис для чтения email
func CreateEmailReader() IEmailReader {
	return &EmailReader{
		imapHost:     os.Getenv("IMAP_HOST"),
		imapPort:     os.Getenv("IMAP_PORT"),
		imapUser:     os.Getenv("IMAP_USER"),
		imapPassword: os.Getenv("IMAP_PASSWORD"),
	}
}

// GetLastMessage получает последнее сообщение из почтового ящика и возвращает его содержимое
func (r *EmailReader) GetLastMessage() (string, error) {
	if r.imapHost == "" {
		return "", fmt.Errorf("IMAP_HOST is not set")
	}
	if r.imapPort == "" {
		return "", fmt.Errorf("IMAP_PORT is not set")
	}
	if r.imapUser == "" {
		return "", fmt.Errorf("IMAP_USER is not set")
	}
	if r.imapPassword == "" {
		return "", fmt.Errorf("IMAP_PASSWORD is not set")
	}

	addr := r.imapHost + ":" + r.imapPort

	// Парсим порт
	port, err := strconv.Atoi(r.imapPort)
	if err != nil {
		return "", fmt.Errorf("invalid IMAP_PORT: %v", err)
	}

	// Подключаемся к IMAP серверу
	var conn net.Conn
	if port == 993 {
		// SSL соединение
		tlsConfig := &tls.Config{
			ServerName:         r.imapHost,
			InsecureSkipVerify: false,
		}
		conn, err = tls.DialWithDialer(
			&net.Dialer{Timeout: 10 * time.Second},
			"tcp",
			addr,
			tlsConfig,
		)
	} else {
		// Обычное соединение
		conn, err = net.DialTimeout("tcp", addr, 10*time.Second)
	}

	if err != nil {
		return "", fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer conn.Close()

	// Устанавливаем таймаут на чтение
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Читаем приветствие сервера
	buffer := make([]byte, 4096)
	_, err = conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read IMAP greeting: %w", err)
	}

	// Отправляем команду LOGIN
	loginCmd := fmt.Sprintf("a001 LOGIN %s %s\r\n", r.imapUser, r.imapPassword)
	_, err = conn.Write([]byte(loginCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send LOGIN command: %w", err)
	}

	// Читаем ответ с таймаутом
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	response, err := r.readIMAPResponse(conn, "a001")
	if err != nil {
		return "", fmt.Errorf("failed to read LOGIN response: %w", err)
	}

	if !strings.Contains(response, "a001 OK") {
		return "", fmt.Errorf("IMAP login failed: %s", response)
	}

	// Выбираем папку INBOX
	selectCmd := "a002 SELECT INBOX\r\n"
	_, err = conn.Write([]byte(selectCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send SELECT command: %w", err)
	}

	// Читаем ответ SELECT
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = r.readIMAPResponse(conn, "a002")
	if err != nil {
		return "", fmt.Errorf("failed to read SELECT response: %w", err)
	}

	// Получаем последнее сообщение (используем UID SEARCH для поиска последнего UID)
	searchCmd := "a003 UID SEARCH ALL\r\n"
	_, err = conn.Write([]byte(searchCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send SEARCH command: %w", err)
	}

	// Читаем ответ SEARCH
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	searchResponse, err := r.readIMAPResponse(conn, "a003")
	if err != nil {
		return "", fmt.Errorf("failed to read SEARCH response: %w", err)
	}

	// Парсим последний UID из ответа (формат: * SEARCH 1 2 3 ...)
	lastUID := "1" // По умолчанию
	if strings.Contains(searchResponse, "SEARCH") {
		// Ищем строку с SEARCH
		lines := strings.Split(searchResponse, "\r\n")
		for _, line := range lines {
			if strings.Contains(line, "SEARCH") && strings.HasPrefix(strings.TrimSpace(line), "*") {
				// Формат: * SEARCH 1 2 3 ...
				parts := strings.Fields(line)
				for i, part := range parts {
					if part == "SEARCH" && i+1 < len(parts) {
						// Берем все UIDs после SEARCH
						uids := parts[i+1:]
						if len(uids) > 0 {
							// Берем последний UID
							lastUID = strings.TrimSpace(uids[len(uids)-1])
							// Убираем возможные символы в конце
							lastUID = strings.Trim(lastUID, " \r\n\t")
						}
						break
					}
				}
				break
			}
		}
	}

	// Проверяем, что UID валидный
	if lastUID == "" || lastUID == "0" {
		return "", fmt.Errorf("no messages found in mailbox")
	}

	// Получаем тело последнего сообщения используя UID FETCH
	fetchCmd := fmt.Sprintf("a004 UID FETCH %s (BODY[TEXT])\r\n", lastUID)
	_, err = conn.Write([]byte(fetchCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send FETCH command: %w", err)
	}

	// Читаем ответ FETCH с увеличенным таймаутом (сообщение может быть большим)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	fetchResponse, err := r.readIMAPResponse(conn, "a004")
	if err != nil {
		return "", fmt.Errorf("failed to read FETCH response: %w", err)
	}

	// Парсим тело сообщения из ответа
	body := fetchResponse

	// Ищем начало тела сообщения (после BODY[TEXT] {размер})
	startIdx := strings.Index(body, "{")
	if startIdx == -1 {
		return "", fmt.Errorf("failed to parse message body: %s", body[:min(200, len(body))])
	}

	endIdx := strings.Index(body[startIdx:], "}")
	if endIdx == -1 {
		return "", fmt.Errorf("failed to parse message body size")
	}

	sizeStr := body[startIdx+1 : startIdx+endIdx]
	size, err := strconv.Atoi(strings.TrimSpace(sizeStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse message size: %w", err)
	}

	// Ищем начало текста сообщения (после \r\n после размера)
	bodyStart := strings.Index(body[startIdx+endIdx:], "\r\n")
	if bodyStart == -1 {
		return "", fmt.Errorf("failed to find message body start")
	}
	bodyStart = startIdx + endIdx + bodyStart + 2

	// Извлекаем текст сообщения
	if bodyStart+size > len(body) {
		// Если сообщение не полностью в буфере, читаем остаток
		remaining := size - (len(body) - bodyStart)
		messageText := body[bodyStart:]

		if remaining > 0 {
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			remainingBuf := make([]byte, remaining)
			n, err := conn.Read(remainingBuf)
			if err != nil {
				return "", fmt.Errorf("failed to read remaining message: %w", err)
			}
			messageText += string(remainingBuf[:n])
		}

		// Убираем лишние символы в конце
		return strings.TrimSpace(messageText), nil
	}

	messageText := body[bodyStart : bodyStart+size]

	// Убираем лишние символы в конце
	messageText = strings.TrimSpace(messageText)

	// Закрываем соединение
	logoutCmd := "a005 LOGOUT\r\n"
	conn.Write([]byte(logoutCmd))

	return messageText, nil
}

// readIMAPResponse читает полный ответ IMAP команды до завершения
func (r *EmailReader) readIMAPResponse(conn net.Conn, tag string) (string, error) {
	var response strings.Builder
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if response.Len() > 0 {
				return response.String(), nil
			}
			return "", fmt.Errorf("failed to read IMAP response: %w", err)
		}

		response.Write(buffer[:n])
		responseStr := response.String()

		// Проверяем, завершилась ли команда (ищем tag OK или tag NO или tag BAD)
		if strings.Contains(responseStr, tag+" OK") ||
			strings.Contains(responseStr, tag+" NO") ||
			strings.Contains(responseStr, tag+" BAD") {
			break
		}

		// Проверяем таймаут
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			break
		}
	}

	return response.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
