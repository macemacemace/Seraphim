export type Photo = {
  id: string;
  url: string;
  position: number;
};




export type Listing = {
  id: string;
  title: string;
  price: number;
  sellerId: string;
  type: string;
  description: string;
  lat: number;
  lng: number;
  photos?: Photo[];
};
