import React from 'react';
import { Typography, Card } from 'antd';

const { Title, Paragraph } = Typography;

const Home: React.FC = () => {
  return (
    <div style={{ padding: '24px' }}>
      <Card>
        <Title level={2}>ПРОЕКТ ПО БАЗАМ ДАННЫХ</Title>
        <Paragraph>
          Приложение предназначено для взаимодействия между заказчиками, репетиторами и сотрудниками. Основные функции включают размещение заказов, подбор репетиторов, ведение чатов, управление учетными записями и отделами, а также систему отзывов. Система позволяет автоматизировать процесс поиска репетиторов и взаимодействия между пользователями, обеспечивая удобство.
        </Paragraph>
      </Card>
    </div>
  );
};

export default Home; 