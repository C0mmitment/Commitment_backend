import rateLimit from "express-rate-limit";

const strictLimiter = rateLimit({
  windowMs: 1 * 60 * 1000,
  max: 20,
});

const normalLimiter = rateLimit({
  windowMs: 1 * 60 * 1000,
  max: 100, 
});

export {
    strictLimiter,
    normalLimiter
}