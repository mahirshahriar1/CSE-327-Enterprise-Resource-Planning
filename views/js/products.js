// product data
const products = [
  {
    name: "Product 1",
    stock: 50,
    icon: "inventory",
    color: "text-blue", // Assuming CSS for this exists
  },
  {
    name: "Product 2",
    stock: 100,
    icon: "inventory",
    color: "text-orange", // Assuming CSS for this exists
  },
  {
    name: "Product 3",
    stock: 11,
    icon: "inventory",
    color: "text-green", // Assuming CSS for this exists
  },
  {
    name: "Product 4",
    stock: 0,
    icon: "inventory",
    color: "text-red", // Assuming CSS for this exists
  },
];

// Function to render products on the page
function renderProducts() {
  const productContainer = document.querySelector(".main-cards");
  productContainer.innerHTML = ""; // Clear any existing content
  console.log("Products array: ", products);
  products.forEach((product) => {
    const productCard = `
              <div class="card">
                <div class="card-inner">
                  <p class="text-primary">${product.name.toUpperCase()}</p>
                  <span class="material-icons-outlined ${product.color}">${
      product.icon
    }</span>
                </div>
                <span class="text-primary font-weight-bold">Stock: ${
                  product.stock
                }</span>
              </div>
            `;
    productContainer.innerHTML += productCard;
  });
}

renderProducts();
