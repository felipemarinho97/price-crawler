import logo from './logo.svg';
import './App.css';
import { Layout } from "antd";
import Page from './page/Page';
import Products from './page/Products';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import Chart from './page/Chart';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Products />,
  },
  {
    path: "/products/:name",
    element: <Chart />,
  },
]);

function App() {
  return (
    <Layout>
      <Layout.Header style={{ color: 'white' }} title='Price Crawler Viewer'>
        Price Crawler Viewer
      </Layout.Header>
      <Layout.Content>
        <div className="app-container" style={{ margin: '1vh 6vw' }}>
          <RouterProvider router={router} />
        </div>
      </Layout.Content>
      <Layout.Footer>
        Price Crawler Viewer ©2021 Created by Price Crawler Team
      </Layout.Footer>
    </Layout>
  );
}

export default App;
