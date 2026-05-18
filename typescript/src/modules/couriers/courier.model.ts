export type CourierStatus = "active" | "inactive" | "suspended";

export interface CourierResponse {
  id: string;
  name: string;
  email: string | null;
  phone: string | null;
  level: number;
  vehicle_type: string | null;
  license_plate: string | null;
  status: string;
  registered_at: Date;
  created_at: Date;
  updated_at: Date;
  deleted_at: Date | null;
}

export interface CourierPayload {
  name: string;
  email?: string | null;
  phone?: string | null;
  level: number;
  vehicle_type?: string | null;
  license_plate?: string | null;
  status?: CourierStatus | null;
  registered_at?: string | null;
}

export interface CourierQuery {
  page: number;
  per_page: number;
  sort: "name" | "-name" | "registered_at" | "-registered_at" | "created_at" | "-created_at";
  search?: string;
  levels: number[];
}
