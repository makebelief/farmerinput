// Initialize Lucide icons
document.addEventListener('DOMContentLoaded', () => {
    lucide.createIcons();
    fetchProducts();
});

// Cart functionality
let cart = [];
const cartCount = document.getElementById('cart-count');

// Product data fetching
async function fetchProducts() {
    try {
        const response = await fetch('http://localhost:8080/api/products');
        const products = await response.json();
        displayProducts(products);
    } catch (error) {
        console.error('Error fetching products:', error);
        showError('Failed to load products. Please try again later.');
    }
}

// Display products in the grid
function displayProducts(products) {
    const container = document.getElementById('product-container');
    container.innerHTML = '';

    products.forEach(product => {
        const productCard = document.createElement('div');
        productCard.className = 'product-card';
        productCard.innerHTML = `
            <div class="product-image">
                <img src="${product.imageUrl}" alt="${product.name}" 
                     onerror="this.src='/api/placeholder/400/300'">
                ${product.stock < 10 ? '<span class="low-stock">Low Stock</span>' : ''}
            </div>
            <div class="product-info">
                <div class="product-header">
                    <h3>
                        <i data-lucide="package"></i>
                        ${product.name}
                    </h3>
                    <span class="brand">${product.brand}</span>
                </div>
                <p class="description">${product.description}</p>
                <div class="product-meta">
                    <p class="product-price">
                        <i data-lucide="tag"></i>
                        KSh ${product.price.toFixed(2)} / ${product.unit}
                    </p>
                    <div class="rating">
                        <span class="stars" style="--rating: ${product.rating}"></span>
                        <span class="reviews">(${product.reviews} reviews)</span>
                    </div>
                </div>
                <div class="product-footer">
                    <p class="product-stock ${product.stock < 20 ? 'low' : ''}">
                        <i data-lucide="layers"></i>
                        ${product.stock} in stock
                    </p>
                    <button class="add-to-cart" onclick="addToCart(${JSON.stringify(product)})">
                        <i data-lucide="shopping-cart"></i>
                        Add to Cart
                    </button>
                </div>
            </div>
        `;
        container.appendChild(productCard);
        lucide.createIcons();
    });
}

// Search products
async function searchProducts(query) {
    try {
        const response = await fetch(`http://localhost:8080/api/products/search?q=${encodeURIComponent(query)}`);
        const products = await response.json();
        displayProducts(products);
    } catch (error) {
        console.error('Error searching products:', error);
        showError('Search failed. Please try again.');
    }
}

// Filter products by category
async function filterByCategory(category) {
    try {
        const response = await fetch(`http://localhost:8080/api/products/category/${encodeURIComponent(category)}`);
        const products = await response.json();
        displayProducts(products);
    } catch (error) {
        console.error('Error filtering products:', error);
        showError('Failed to filter products. Please try again.');
    }
}

// Show error message
function showError(message) {
    const notification = document.createElement('div');
    notification.className = 'notification error';
    notification.innerHTML = `
        <i data-lucide="alert-circle"></i>
        ${message}
    `;
    document.body.appendChild(notification);
    lucide.createIcons();
    
    setTimeout(() => {
        notification.remove();
    }, 5000);
}

// Rest of your existing cart functionality...

// Add to cart functionality
function addToCart(product) {
    const existingItem = cart.find(item => item.id === product.id);
    
    if (existingItem) {
        existingItem.quantity += 1;
    } else {
        cart.push({ ...product, quantity: 1 });
    }
    
    updateCartCount();
    showNotification(`Added ${product.name} to cart`);
}

// Update cart count in the navigation
function updateCartCount() {
    cartCount.textContent = cart.reduce((total, item) => total + item.quantity, 0);
}

// Show notification
function showNotification(message) {
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.innerHTML = `
        <i data-lucide="check-circle"></i>
        ${message}
    `;
    document.body.appendChild(notification);
    lucide.createIcons();
    
    setTimeout(() => {
        notification.remove();
    }, 3000);
}

// Initialize
fetchProducts();