#!/usr/bin/env node
/**
 * Auto-generate typescript/index.ts from generated.ts
 *
 * This script:
 * 1. Reads typescript/generated.ts
 * 2. Extracts all schema names from components["schemas"]
 * 3. Generates convenience type exports in index.ts
 *
 * Run: node scripts/generate-index.js
 */

const fs = require('fs');
const path = require('path');

const GENERATED_FILE = path.join(__dirname, '../src/generated.ts');
const OUTPUT_FILE = path.join(__dirname, '../src/index.ts');

// Read generated.ts
const generatedContent = fs.readFileSync(GENERATED_FILE, 'utf-8');

// Extract schema names from components["schemas"]
// Match objects: SchemaName: { ... } or SchemaName: components["schemas"]["X"] & { ... }
// And enums: SchemaName: "val1" | "val2";
const objectSchemas = generatedContent.matchAll(/^\s{8}(\w+):\s*(?:components\["schemas"\]\["\w+"\]\s*&\s*)?{/gm);
const enumSchemas = generatedContent.matchAll(/^\s{8}(\w+):\s*"/gm);

const schemaNames = [
  ...Array.from(objectSchemas, m => m[1]),
  ...Array.from(enumSchemas, m => m[1])
];

console.log(`Found ${schemaNames.length} schemas (objects + enums)`);

// Categorize schemas
const categories = {
  entity: ['User', 'Manga', 'UserProgress', 'UserProgressWithManga'],
  enum: ['MangaStatus', 'ReadingStatus', 'TCPMessageType', 'UDPMessageType', 'ChatMessageType'],
  request: ['UserRegisterRequest', 'UserLoginRequest', 'LibraryAddRequest', 'ProgressUpdateRequest'],
  response: [
    'APIResponse', 'APIError', 'Meta',
    'AuthResponse', 'UserProfileResponse',
    'MangaListResponse', 'MangaDetailResponse',
    'LibraryResponse', 'ProgressResponse'
  ],
  tcp: [
    'TCPMessage', 'TCPAuthMessage', 'TCPAuthSuccessMessage', 'TCPAuthFailedMessage',
    'TCPProgressMessage', 'TCPProgressBroadcast', 'TCPErrorMessage', 'TCPPingMessage', 'TCPPongMessage'
  ],
  udp: [
    'UDPMessage', 'UDPRegisterMessage', 'UDPRegisterSuccessMessage', 'UDPRegisterFailedMessage',
    'UDPUnregisterMessage', 'UDPNotification', 'UDPPingMessage', 'UDPPongMessage', 'UDPErrorMessage'
  ],
  websocket: ['ChatMessage']
};

// Generate index.ts content
const header = `/**
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
 * Generated: ${new Date().toISOString()}
 */

// Re-export all generated types
export * from './generated';

// ============================================
// Additional Utility Types
// ============================================

import type { components, paths } from './generated';

// Convenient type aliases
export type Schemas = components['schemas'];
`;

const sections = [];

// Entity types
sections.push(`// ============================================
// Entity Types
// ============================================

${categories.entity.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// Enum types
sections.push(`// ============================================
// Enum Types
// ============================================

${categories.enum.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// Request types
sections.push(`// ============================================
// Request Types
// ============================================

${categories.request.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// Response types
sections.push(`// ============================================
// Response Types
// ============================================

${categories.response.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// TCP message types
sections.push(`// ============================================
// TCP Message Types
// ============================================

${categories.tcp.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// UDP message types
sections.push(`// ============================================
// UDP Message Types
// ============================================

${categories.udp.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// WebSocket message types
sections.push(`// ============================================
// WebSocket Message Types
// ============================================

${categories.websocket.filter(s => schemaNames.includes(s)).map(s => `export type ${s} = Schemas['${s}'];`).join('\n')}`);

// Path operation types
const footer = `
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
`;

const content = header + sections.join('\n\n') + footer;

// Write index.ts
fs.writeFileSync(OUTPUT_FILE, content, 'utf-8');

console.log('✅ Generated typescript/index.ts');
console.log(`   - Entity types: ${categories.entity.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - Enum types: ${categories.enum.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - Request types: ${categories.request.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - Response types: ${categories.response.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - TCP message types: ${categories.tcp.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - UDP message types: ${categories.udp.filter(s => schemaNames.includes(s)).length}`);
console.log(`   - WebSocket types: ${categories.websocket.filter(s => schemaNames.includes(s)).length}`);
