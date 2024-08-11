import React from 'react';
import useSWR from 'swr';
import { useParams, useNavigate } from 'react-router';
import { Line } from '@ant-design/charts';
import { format } from 'fecha';
import { Button, Card, Spin } from 'antd';

const fetcher = (...args) => fetch(...args).then((res) => res.json());

const startTime = new Date(2024, 1, 1);
const endTime = new Date();
// get product name from route params
const Chart = () => {
    const params = useParams();
    const {
        data,
        error,
        isValidating,
      } = useSWR(`http://localhost:9081/datapoints?name=${params.name}&start=${startTime.toISOString()}&end=${endTime.toISOString()}`, fetcher);
    const navigate = useNavigate();
    if (error) {
        return <Card title={params.name}>Failed to load data. {error.message}</Card>;
    }

    if (isValidating) {
        return <Card title={params.name}><Spin size='large' /></Card>;
    }

      const props = {
        data,
        xField: (d) => new Date(d.timestamp),
        yField: 'value',
        slider: {
            x: { labelFormatter: (d) => format(d, 'YYYY/M/D') },
            y: { labelFormatter: '~s' },
        },
        axis: {
            y: { title: '↑ Product Price (R$)' },
        },
        label: {
            y: (d) => ({ content: `R$ ${d.value}` }),
        },
      };

      return (
        <Card title={params.name}>
            <Line {...props} />
            <Button
             type="primary" onClick={() => navigate(`/`)}
                style={{ marginTop: '1vh' }}  
            >
                Back
            </Button>
        </Card>
      );
}

export default Chart;