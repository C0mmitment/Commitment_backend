import cors from 'cors';

const corsOptions = {
  origin: ['http://localhost:8081','http://localhost:5173','http://localhost:5174','https://c8t.esrj86.org/'], // フロントエンドのURLに合わせて変更（例: React）
  credentials: true,              // Cookieを許可
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: [
    'Content-Type',  // multipart/form-data含む
    'Authorization',
    // 'X-Requested-With',
  ],
};

export default cors(corsOptions);