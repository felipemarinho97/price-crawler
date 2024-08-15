import React from 'react';
import { Button, Card, List, Space } from 'antd';
import useSWR from 'swr';
import { DollarCircleOutlined, FallOutlined, RiseOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { DATA_BUCKET_URL } from '../constants';
import './Products.css';

const fetcher = (...args) => fetch(...args).then((res) => res.json());

const IconText = ({ icon, text }) => (
    <Space>
      {React.createElement(icon)}
      {text}
    </Space>
  );

const Actions = ({ item }) => {
    let actions = [];
    const cheapPriceEver = item.lastValue <= item.minValue;
    const cheapPriceNow = item.lastValue <= item.avgValue;
    const expensivePriceNow = item.lastValue >= item.avgValue;
    const nonVariablePrice = item.minValue === item.maxValue;

    if (cheapPriceEver && cheapPriceNow && !nonVariablePrice) {
        actions.push(<IconText icon={FallOutlined} text={<span className='cheap'>The cheapest price ever!</span>} key="cheapEverNow"/>);
    }
    
    if (expensivePriceNow && !nonVariablePrice) {
        actions.push(<IconText icon={RiseOutlined} text={<span className='expensive'>Not buy! The price is too high!</span>} key="expensiveNow"/>);
    }

    actions.push(<IconText icon={DollarCircleOutlined} text={<span className='price'>Price: R$ {item.lastValue.toFixed(2)}</span>} key="price"/>);

    return actions.reduce((prev, curr) => [prev, '  ', curr]);
}

const Products = () => {
    const {
        data,
        error,
        isValidating,
      } = useSWR(`${DATA_BUCKET_URL}/datapoints/name`, fetcher);
    const navigate = useNavigate();
    
    return (
        <Card title="Products">
            <List
            itemLayout="horizontal"
            dataSource={data}
            renderItem={item => (
                <List.Item
                    actions={[
                        <Actions item={item} />,
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
