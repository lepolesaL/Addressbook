export class Contact{
    name: String;
    email: String;
    phonenumber: String;
    address: Address;
}

interface Address {
    street: String;
    city: String;
    country: String;
  }

export interface ApiError {
      ErrorCode: String;
      Message: String;
}