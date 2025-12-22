/**
 * MangaHub API TypeScript Types
 *
 * THIS FILE IS AUTO-GENERATED - DO NOT EDIT MANUALLY
 *
 * To regenerate:
 *   yarn workspace @mangahub/types generate
 *
 * Or from project root:
 *   make generate-types
 *
 * Generated: 2025-12-22T18:39:36.189Z
 */

// Re-export all generated types
export * from './generated';

// ============================================
// Additional Utility Types
// ============================================

import type { components, paths } from './generated';

// Convenient type aliases
export type Schemas = components['schemas'];
// ============================================
// Entity Types
// ============================================

export type User = Schemas['User'];
export type Manga = Schemas['Manga'];
export type UserProgress = Schemas['UserProgress'];
export type UserProgressWithManga = Schemas['UserProgressWithManga'];

// ============================================
// Enum Types
// ============================================

export type MangaStatus = Schemas['MangaStatus'];
export type ReadingStatus = Schemas['ReadingStatus'];
export type TCPMessageType = Schemas['TCPMessageType'];
export type UDPMessageType = Schemas['UDPMessageType'];
export type ChatMessageType = Schemas['ChatMessageType'];

// ============================================
// Request Types
// ============================================

export type UserRegisterRequest = Schemas['UserRegisterRequest'];
export type UserLoginRequest = Schemas['UserLoginRequest'];
export type LibraryAddRequest = Schemas['LibraryAddRequest'];
export type ProgressUpdateRequest = Schemas['ProgressUpdateRequest'];

// ============================================
// Response Types
// ============================================

export type APIResponse = Schemas['APIResponse'];
export type APIError = Schemas['APIError'];
export type Meta = Schemas['Meta'];
export type AuthResponse = Schemas['AuthResponse'];
export type UserProfileResponse = Schemas['UserProfileResponse'];
export type MangaListResponse = Schemas['MangaListResponse'];
export type MangaDetailResponse = Schemas['MangaDetailResponse'];
export type LibraryResponse = Schemas['LibraryResponse'];
export type ProgressResponse = Schemas['ProgressResponse'];

// ============================================
// TCP Message Types
// ============================================

export type TCPMessage = Schemas['TCPMessage'];
export type TCPAuthMessage = Schemas['TCPAuthMessage'];
export type TCPAuthSuccessMessage = Schemas['TCPAuthSuccessMessage'];
export type TCPAuthFailedMessage = Schemas['TCPAuthFailedMessage'];
export type TCPProgressMessage = Schemas['TCPProgressMessage'];
export type TCPProgressBroadcast = Schemas['TCPProgressBroadcast'];
export type TCPErrorMessage = Schemas['TCPErrorMessage'];

// ============================================
// UDP Message Types
// ============================================

export type UDPMessage = Schemas['UDPMessage'];
export type UDPRegisterMessage = Schemas['UDPRegisterMessage'];
export type UDPRegisterSuccessMessage = Schemas['UDPRegisterSuccessMessage'];
export type UDPRegisterFailedMessage = Schemas['UDPRegisterFailedMessage'];
export type UDPUnregisterMessage = Schemas['UDPUnregisterMessage'];
export type UDPNotification = Schemas['UDPNotification'];
export type UDPPingMessage = Schemas['UDPPingMessage'];
export type UDPPongMessage = Schemas['UDPPongMessage'];
export type UDPErrorMessage = Schemas['UDPErrorMessage'];

// ============================================
// WebSocket Message Types
// ============================================

export type ChatMessage = Schemas['ChatMessage'];
// ============================================
// Path Operation Types (for API client typing)
// ============================================

export type Paths = paths;

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
