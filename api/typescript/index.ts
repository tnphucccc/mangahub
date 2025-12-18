/**
 * MangaHub API TypeScript Types
 * 
 * THIS FILE IS AUTO-GENERATED FROM openapi.yaml
 * DO NOT EDIT MANUALLY
 * 
 * To regenerate:
 *   cd api && npm run generate
 * 
 * Or from project root:
 *   make generate-types
 */

// Re-export all generated types
export * from './generated';

// ============================================
// Additional Utility Types
// ============================================

import type { components, paths } from './generated';

// Convenient type aliases
export type Schemas = components['schemas'];

// Entity types
export type User = Schemas['User'];
export type Manga = Schemas['Manga'];
export type UserProgress = Schemas['UserProgress'];
export type UserProgressWithManga = Schemas['UserProgressWithManga'];

// Enum types
export type MangaStatus = Schemas['MangaStatus'];
export type ReadingStatus = Schemas['ReadingStatus'];
export type TCPMessageType = Schemas['TCPMessageType'];
export type ChatMessageType = Schemas['ChatMessageType'];

// Request types
export type UserRegisterRequest = Schemas['UserRegisterRequest'];
export type UserLoginRequest = Schemas['UserLoginRequest'];
export type LibraryAddRequest = Schemas['LibraryAddRequest'];
export type ProgressUpdateRequest = Schemas['ProgressUpdateRequest'];

// Response types
export type APIResponse = Schemas['APIResponse'];
export type APIError = Schemas['APIError'];
export type Meta = Schemas['Meta'];
export type AuthResponse = Schemas['AuthResponse'];
export type MangaListResponse = Schemas['MangaListResponse'];
export type MangaDetailResponse = Schemas['MangaDetailResponse'];
export type LibraryResponse = Schemas['LibraryResponse'];
export type ProgressResponse = Schemas['ProgressResponse'];

// TCP/WebSocket message types
export type TCPMessage = Schemas['TCPMessage'];
export type TCPAuthMessage = Schemas['TCPAuthMessage'];
export type TCPAuthSuccessMessage = Schemas['TCPAuthSuccessMessage'];
export type TCPAuthFailedMessage = Schemas['TCPAuthFailedMessage'];
export type TCPProgressMessage = Schemas['TCPProgressMessage'];
export type TCPProgressBroadcast = Schemas['TCPProgressBroadcast'];
export type TCPErrorMessage = Schemas['TCPErrorMessage'];
export type ChatMessage = Schemas['ChatMessage'];

// Path operation types (for API client typing)
export type Paths = paths;

// ============================================
// API Client Helper Types
// ============================================

/** Extract successful response data type from an endpoint */
export type SuccessResponse<T extends keyof paths, M extends keyof paths[T]> = 
  paths[T][M] extends { responses: { 200: { content: { 'application/json': infer R } } } }
    ? R
    : paths[T][M] extends { responses: { 201: { content: { 'application/json': infer R } } } }
    ? R
    : never;

/** Extract request body type from an endpoint */
export type RequestBody<T extends keyof paths, M extends keyof paths[T]> =
  paths[T][M] extends { requestBody: { content: { 'application/json': infer R } } }
    ? R
    : never;

/** Extract query parameters type from an endpoint */
export type QueryParams<T extends keyof paths, M extends keyof paths[T]> =
  paths[T][M] extends { parameters: { query?: infer Q } }
    ? Q
    : never;
