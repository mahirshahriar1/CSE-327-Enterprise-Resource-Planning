// Updated product data aligning with the Product model
const products = [
  {
    id: 1,
    name: "Summer T-Shirt",
    brand: "Urban Threads",
    season: "Summer",
    price: "$15.99",
    stock: 50,
    icon: "inventory",
    color: "text-blue", // Assuming CSS for this exists
  },
  {
    id: 2,
    name: "Winter Jacket",
    brand: "Arctic Wear",
    season: "Winter",
    price: "$89.99",
    stock: 100,
    icon: "inventory",
    color: "text-orange",
  },
  {
    id: 3,
    name: "Casual Jeans",
    brand: "Denim Co.",
    season: "All",
    price: "$49.99",
    stock: 11,
    icon: "inventory",
    color: "text-green",
  },
  {
    id: 4,
    name: "Raincoat",
    brand: "DrizzleSafe",
    season: "Monsoon",
    price: "$39.99",
    stock: 0,
    icon: "inventory",
    color: "text-red",
  },
];

// Function to render products on the page
function renderProducts() {
  const productContainer = document.querySelector(".main-cards");
  productContainer.innerHTML = ""; // Clear any existing content

  products.forEach((product) => {
    const productCard = `
              <div class="card">
                <div class="card-inner">
                  <p class="text-primary">${product.name.toUpperCase()}</p>
                  <span class="material-icons-outlined ${product.color}">${product.icon}</span>
                </div>
                <span class="text-primary font-weight-bold">Brand: ${product.brand}</span>
                <span class="text-primary font-weight-bold">Season: ${product.season}</span>
                <span class="text-primary font-weight-bold">Price: ${product.price}</span>
                <span class="text-primary font-weight-bold">Stock: ${product.stock}</span>
              </div>
            `;
    productContainer.innerHTML += productCard;
  });
}

renderProducts();