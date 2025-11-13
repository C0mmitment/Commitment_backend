import cors from 'cors';

const corsOptions = {
  origin: 'http://localhost:8081', // フロントエンドのURLに合わせて変更（例: React）
  credentials: true,              // Cookieを許可
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: [
    'Content-Type',  // multipart/form-data含む
    'Authorization',
    // 'X-Requested-With',
  ],
};

export default cors(corsOptions);