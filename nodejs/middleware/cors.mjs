import cors from 'cors';

const corsOptions = {
  origin: ['http://localhost:8081','http://localhost:5173','http://localhost:5174','https://totteme.esrj86.org/'], 
  credentials: true,             
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: [
    'Content-Type', 
    'Authorization',
  ],
};

export default cors(corsOptions);