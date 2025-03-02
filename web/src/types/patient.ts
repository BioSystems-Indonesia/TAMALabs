import type { Specimen } from "./specimen";

export interface Patient {
  id: number;
  first_name: string;
  last_name: string;
  birthdate: string;
  sex: string;
  phone_number: string;
  location: string;
  address: string;
  created_at: string;
  updated_at: string;
  specimen_list: Specimen[] | null;
}

