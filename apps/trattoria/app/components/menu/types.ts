export interface IMenuItem {
  name: string;
  description: string;
  price: number;
  image: {
    fields: {
      file: {
        url: string;
      };
    };
  };
}
