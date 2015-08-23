CREATE TABLE items (
	  id serial ,
	  Title varchar (255) NOT NULL,
	  Date  date  NOT NULL,
		Seller  varchar (255) NOT NULL,
		Price   money,
		Image varchar (255) NOT NULL,
		Description text
);
