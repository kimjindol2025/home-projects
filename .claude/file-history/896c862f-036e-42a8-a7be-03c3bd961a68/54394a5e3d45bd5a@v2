/**
 * Input Validation Middleware using Joi
 */

import Joi from 'joi';
import { Request, Response, NextFunction } from 'express';

// OAuth Callback Validation Schema
export const oauthCallbackSchema = Joi.object({
  code: Joi.string().alphanum().min(10).max(500).required(),
  state: Joi.string().max(200).optional(),
});

// Create Post Validation Schema
export const createPostSchema = Joi.object({
  title: Joi.string().min(1).max(200).required(),
  content: Joi.string().min(1).max(50000).required(),
});

// Token Exchange Validation Schema
export const tokenExchangeSchema = Joi.object({
  code: Joi.string().alphanum().min(10).max(500).required(),
});

/**
 * Generic Validation Middleware Factory
 * Validates request data against a Joi schema
 *
 * @param schema - Joi schema to validate against
 * @param source - Request property to validate ('query', 'body', 'params', 'headers')
 */
export const validate = (
  schema: Joi.ObjectSchema,
  source: 'query' | 'body' | 'params' | 'headers' = 'body'
) => {
  return (req: Request, res: Response, next: NextFunction) => {
    const { error, value } = schema.validate(req[source], {
      abortEarly: false,
      stripUnknown: true,
    });

    if (error) {
      const details = error.details.map((detail) => ({
        field: detail.path.join('.'),
        message: detail.message,
      }));

      return res.status(400).json({
        error: 'Validation failed',
        details,
        timestamp: new Date().toISOString(),
      });
    }

    // Replace original data with validated/cleaned data
    req[source] = value;
    next();
  };
};
