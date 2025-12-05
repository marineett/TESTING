package service_logic

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

// IEmailSender интерфейс для отправки email
type IEmailSender interface {
	SendEmail(to, subject, body string) error
}

// EmailSender сервис для отправки email через SMTP
type EmailSender struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// CreateEmailSender создает сервис для отправки email
func CreateEmailSender() IEmailSender {
	return &EmailSender{
		smtpHost:     os.Getenv("SENDER_SMTP_HOST"),
		smtpPort:     os.Getenv("SENDER_SMTP_PORT"),
		smtpUser:     os.Getenv("SENDER_SMTP_USER"),
		smtpPassword: os.Getenv("SENDER_SMTP_PASSWORD"),
		fromEmail:    os.Getenv("SENDER_SMTP_FROM_EMAIL"),
		fromName:     os.Getenv("SENDER_SMTP_FROM_NAME"),
	}
}

func (s *EmailSender) SendEmail(to string, subject string, body string) error {
	if s.smtpHost == "" {
		return fmt.Errorf("SMTP_HOST is not set")
	}
	if s.smtpPort == "" {
		return fmt.Errorf("SMTP_PORT is not set")
	}
	if s.fromEmail == "" {
		return fmt.Errorf("SMTP_FROM_EMAIL is not set")
	}

	// Парсим порт
	port, err := strconv.Atoi(s.smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	// Формируем сообщение
	msg := []byte("From: " + s.fromEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := s.smtpHost + ":" + s.smtpPort

	// Для порта 465 используем SSL напрямую, для 587 - STARTTLS
	if port == 465 {
		return s.sendWithSSL(addr, to, msg)
	}
	return s.sendWithSTARTTLS(addr, to, msg)
}

// sendWithSSL отправляет email через SSL (порт 465)
func (s *EmailSender) sendWithSSL(addr string, to string, msg []byte) error {
	// Создаем TLS конфигурацию
	tlsConfig := &tls.Config{
		ServerName:         s.smtpHost,
		InsecureSkipVerify: false,
	}

	// Устанавливаем соединение с таймаутом
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 10 * time.Second},
		"tcp",
		addr,
		tlsConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Создаем SMTP клиент
	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Авторизация
	if s.smtpUser != "" && s.smtpPassword != "" {
		auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Отправка
	if err := client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS отправляет email через STARTTLS (порт 587)
func (s *EmailSender) sendWithSTARTTLS(addr string, to string, msg []byte) error {
	// Устанавливаем соединение с таймаутом
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Создаем SMTP клиент
	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Проверяем поддержку STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         s.smtpHost,
			InsecureSkipVerify: false,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Авторизация
	if s.smtpUser != "" && s.smtpPassword != "" {
		auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Отправка
	if err := client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

