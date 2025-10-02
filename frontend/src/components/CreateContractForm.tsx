import React, { useState } from 'react';
import { Form, Input, Button, Card, message, InputNumber, Select } from 'antd';
import { ContractInitInfo } from '../types/contract';
import { API_ENDPOINTS } from '../config';

interface CreateContractFormProps {
    clientId: number;
    onContractCreated?: () => void;
}

const CreateContractForm: React.FC<CreateContractFormProps> = ({ clientId, onContractCreated }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState<boolean>(false);

    const CONTRACT_CATEGORIES = [
        { value: 1, label: 'Перевод' },
        { value: 2, label: 'Написание' },
        { value: 3, label: 'Дизайн' },
        { value: 4, label: 'Программирование' },
        { value: 5, label: 'Другое' },
    ];
    const CONTRACT_SUBCATEGORIES = [
        { value: 1, label: 'Репетиторство' },
        { value: 2, label: 'Перевод' },
        { value: 3, label: 'Написание' },
        { value: 4, label: 'Дизайн' },
        { value: 5, label: 'Программирование' },
        { value: 6, label: 'Другое' },
    ];

    const onFinish = async (values: any): Promise<void> => {
        setLoading(true);
        try {
            let isoDate = values.start_date;
            if (isoDate && !isoDate.includes('T')) {
                isoDate = isoDate + 'T00:00:00Z';
            }

            const contractData: ContractInitInfo = {
                ...values,
                client_id: clientId,
                start_date: isoDate,
                price: Number(values.price),
                commission: Number(values.commission),
                duration: Number(values.duration),
                contract_category: Number(values.contract_category),
                contract_subcategories: Array.isArray(values.contract_subcategories)
                    ? values.contract_subcategories.map(Number)
                    : [],
            };

            const response = await fetch(API_ENDPOINTS.CLIENT.CREATE_CONTRACT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(contractData),
            });

            if (!response.ok) throw new Error('Failed to create contract');
            await response.json();
            message.success('Контракт успешно создан!');
            form.resetFields();
            if (onContractCreated) onContractCreated();
        } catch (error) {
            message.error('Ошибка при создании контракта');
            console.error('Error:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card title="Создать новый контракт" style={{ marginBottom: 20 }}>
            <Form
                form={form}
                layout="vertical"
                onFinish={onFinish}
            >
                <Form.Item
                    name="title"
                    label="Название контракта"
                    rules={[{ required: true, message: 'Пожалуйста, введите название контракта!' }]}
                >
                    <Input id="title" placeholder="Введите название контракта" />
                </Form.Item>

                <Form.Item
                    name="description"
                    label="Описание"
                    rules={[{ required: true, message: 'Пожалуйста, введите описание контракта!' }]}
                >
                    <Input.TextArea id="description" rows={4} placeholder="Введите описание контракта" />
                </Form.Item>

                <Form.Item
                    name="price"
                    label="Цена (₽)"
                    rules={[{ required: true, message: 'Пожалуйста, введите цену контракта!' }]}
                >
                    <InputNumber id="price" min={0} style={{ width: '100%' }} placeholder="Введите цену контракта" />
                </Form.Item>

                <Form.Item
                    name="start_date"
                    label="Дата начала"
                    rules={[{ required: true, message: 'Пожалуйста, выберите дату начала!' }]}
                >
                    <Input id="start_date" type="date" placeholder="Выберите дату начала" />
                </Form.Item>

                <Form.Item
                    name="contract_category"
                    label="Категория контракта"
                    rules={[{ required: true, message: 'Выберите категорию!' }]}
                >
                    <Select id="contract_category" style={{ width: '100%' }} placeholder="Выберите категорию">
                        {CONTRACT_CATEGORIES.map(opt => (
                            <Select.Option key={opt.value} value={opt.value}>{opt.label}</Select.Option>
                        ))}
                    </Select>
                </Form.Item>

                <Form.Item
                    name="commission"
                    label="Комиссия (₽)"
                    rules={[{ required: true, message: 'Введите комиссию!' }]}
                >
                    <InputNumber id="commission" min={0} style={{ width: '100%' }} placeholder="Введите комиссию" />
                </Form.Item>

                <Form.Item
                    name="duration"
                    label="Длительность (дней)"
                    rules={[{ required: true, message: 'Введите длительность!' }]}
                >
                    <InputNumber id="duration" min={1} style={{ width: '100%' }} placeholder="Введите длительность" />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        Создать контракт
                    </Button>
                </Form.Item>
            </Form>
        </Card>
    );
};

export default CreateContractForm; 