import React from 'react';
import { Button, Card, List, Space } from 'antd';
import useSWR from 'swr';
import { StockOutlined, FallOutlined, RiseOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const fetcher = (...args) => fetch(...args).then((res) => res.json());

const IconText = ({ icon, text }) => (
    <Space>
      {React.createElement(icon)}
      {text}
    </Space>
  );

const Products = () => {
    const {
        data,
        error,
        isValidating,
      } = useSWR('http://localhost:9081/datapoints/name', fetcher);
    // const data = [
    //     { title: 'Product 1', price: 100 },
    //     { title: 'Product 2', price: 200 },
    //     { title: 'Product 3', price: 300 },
    //     { title: 'Product 4', price: 400 },
    //     { title: 'Product 5', price: 500 },
    // ];
    const navigate = useNavigate();
    
    return (
        <Card title="Products">
            <List
            itemLayout="horizontal"
            dataSource={data}
            renderItem={item => (
                <List.Item
                    actions={[
                        <IconText icon={StockOutlined} text={`R$ ${item.avgValue}`} key="avgValue" />,
                        <IconText icon={FallOutlined} text={`R$ ${item.minValue}`} key="minValue" />,
                        <IconText icon={RiseOutlined} text={`R$ ${item.maxValue}`} key="maxValue" />,
                        <Button type="primary" onClick={() => navigate(`/products/${item.name}`)}>
                            View
                        </Button>,
                    ]}
                >
                <List.Item.Meta
                    title={<a onClick={() => navigate(`/products/${item.name}`)}>{item.name}</a>}
                />
                </List.Item>
            )}
            />
        </Card>
    );
}

export default Products;
