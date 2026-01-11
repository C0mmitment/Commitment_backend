import rateLimit from "express-rate-limit";

const strictLimiter = rateLimit({
  windowMs: 1 * 60 * 1000,
  max: 20, // AI・アップロード用
});

const normalLimiter = rateLimit({
  windowMs: 1 * 60 * 1000,
  max: 100, // 参照系
});

export {
    strictLimiter,
    normalLimiter
}