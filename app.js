// Initialize Lucide icons
document.addEventListener('DOMContentLoaded', () => {
    lucide.createIcons();
    
    // Show loading state
    const container = document.getElementById('product-container');
    container.innerHTML = '<div class="loading">Loading products...</div>';
    
    fetchProducts();
});

// Cart functionality
let cart = [];
const cartCount = document.getElementById('cart-count');

// Product data fetching with loading state
async function fetchProducts() {
    try {
        const container = document.getElementById('product-container');
        
        // Fetch the data
        const response = await fetch('http://localhost:8080/api/products');
        const products = await response.json();
        
        // Clear the container once
        container.innerHTML = '';
        
        // Create a document fragment to batch DOM updates
        const fragment = document.createDocumentFragment();
        
        // Create all product cards
        products.forEach(product => {
            const productCard = createProductCard(product);
            fragment.appendChild(productCard);
        });
        
        // Add all cards to the DOM at once
        container.appendChild(fragment);
        
        // Initialize all Lucide icons at once after adding products
        lucide.createIcons();
        
    } catch (error) {
        console.error('Error fetching products:', error);
        showError('Failed to load products. Please try again later.');
    }
}

// Create product card function
function createProductCard(product) {
    const productCard = document.createElement('div');
    productCard.className = 'product-card';
    
    // Pre-load image
    const img = new Image();
    img.src = product.imageUrl;
    img.onerror = () => img.src = '/api/placeholder/400/300';
    
    productCard.innerHTML = `
        <div class="product-image">
            <img src="${product.imageUrl}" alt="${product.name}" 
                 onerror="this.src='/api/placeholder/400/300'"
                 loading="lazy">
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
                <button class="add-to-cart" onclick="addToCart(${JSON.stringify(product).replace(/"/g, '&quot;')})">
                    <i data-lucide="shopping-cart"></i>
                    Add to Cart
                </button>
            </div>
        </div>
    `;
    
    return productCard;
}

// Add some CSS for loading state
const style = document.createElement('style');
style.textContent = `
    .loading {
        text-align: center;
        padding: 2rem;
        font-size: 1.2rem;
        color: #666;
    }
    
    .product-card {
        opacity: 1;
        transition: opacity 0.3s ease-in-out;
    }
    
    .product-card img {
        opacity: 0;
        transition: opacity 0.3s ease-in-out;
    }
    
    .product-card img.loaded {
        opacity: 1;
    }
`;
document.head.appendChild(style);

// Rest of your existing cart functionality...
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

function updateCartCount() {
    cartCount.textContent = cart.reduce((total, item) => total + item.quantity, 0);
}

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