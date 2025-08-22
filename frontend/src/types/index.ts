export interface Service {
  id: number;
  name: string;
  description: string;
  package: string;
  port: number;
  status: boolean;
  configuration: string;
  output_path: string;
  created_at: string;
  updated_at: string;
  endpoints?: Endpoint[];
}

export interface Endpoint {
  id: number;
  service_id: number;
  method: string;
  path: string;
  handler: string;
  table: string;
  operation: string;
  status: boolean;
  auth: boolean;
  created_at: string;
  updated_at: string;
}

export interface ServiceConfig {
  service: ServiceInfo;
  database: DatabaseConfig;
  tables: TableConfig[];
  endpoints: EndpointConfig[];
}

export interface ServiceInfo {
  name: string;
  description: string;
  package: string;
  port: number;
}

export interface DatabaseConfig {
  type: string;
  host: string;
  port: number;
  database: string;
  username: string;
  password: string;
}

export interface TableConfig {
  name: string;
  fields: FieldConfig[];
}

export interface FieldConfig {
  name: string;
  type: string;
  primary_key?: boolean;
  required?: boolean;
  unique?: boolean;
  default?: any;
  validation?: Record<string, string>;
}

export interface EndpointConfig {
  method: string;
  path: string;
  handler: string;
  table?: string;
  operation: string;
  auth?: boolean;
  middleware?: string[];
  request_body?: any;
  response?: any;
}

export interface GenerationResult {
  success: boolean;
  service_name: string;
  output_path: string;
  files: string[];
  endpoints: string[];
  message: string;
  error?: string;
}