export type SearchType = 'Search by Text and Sound' | 'Search by Text' | 'Search by Sound';

export interface SearchResult {
  allowed: string[];
  notAllowed: string[];
}

export interface FeedbackRequest {
  name: string;
  email: string;
  subject: string;
  message: string;
  source: 'web';
}

export interface SuggestionRequest {
  inputText: string;
  allowed: boolean;
}

export interface CategoriesResult {
  allowed: string[];
  notAllowed: string[];
}

export type CategoryKind = 'Allowed' | 'Not Allowed';
