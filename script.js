// Global variables
let currentStep = 1;
let selectedCurrency = '';
let addressCount = 1;
let btcPrice = 43000; // Simulated BTC price in USD

// Simulated cryptocurrency prices (relative to BTC)
const cryptoPrices = {
    'ETH': 0.06, // 1 BTC = 0.06 ETH ratio (approximate)
    'USDT': 43000, // 1 BTC = 43000 USDT
    'USDC': 43000, // 1 BTC = 43000 USDC
    'ADA': 86000, // 1 BTC = 86000 ADA (approximate)
    'SOL': 400, // 1 BTC = 400 SOL (approximate)
    'MATIC': 48000 // 1 BTC = 48000 MATIC (approximate)
};

// Initialize the app
document.addEventListener('DOMContentLoaded', function() {
    updateBTCAmount();
    setupEventListeners();
    generatePaymentAddress();
});

function setupEventListeners() {
    // BTC amount input
    const btcAmountInput = document.getElementById('btcAmount');
    if (btcAmountInput) {
        btcAmountInput.addEventListener('input', updateBTCAmount);
    }

    // Currency selection
    const currencyOptions = document.querySelectorAll('.currency-option');
    currencyOptions.forEach(option => {
        option.addEventListener('click', selectCurrency);
    });

    // Address input validation
    document.addEventListener('input', function(e) {
        if (e.target.matches('.address-input input[type="text"]')) {
            validateAddresses();
        }
        if (e.target.matches('.percentage')) {
            updatePercentages();
        }
    });
}

function updateBTCAmount() {
    const btcAmount = parseFloat(document.getElementById('btcAmount').value) || 0;
    const usdAmount = btcAmount * btcPrice;
    
    document.getElementById('amountUSD').textContent = `≈ $${usdAmount.toLocaleString()} USD`;
    
    // Update summary if we're on step 4
    if (currentStep >= 4) {
        updateOrderSummary();
    }
}

function selectCurrency(e) {
    // Remove previous selection
    document.querySelectorAll('.currency-option').forEach(opt => {
        opt.classList.remove('selected');
    });
    
    // Add selection to clicked option
    e.currentTarget.classList.add('selected');
    selectedCurrency = e.currentTarget.dataset.currency;
    
    // Enable continue button
    const nextBtn = document.querySelector('[data-step="2"] .btn-next');
    nextBtn.disabled = false;
    
    // Update text in step 3
    const currencySpan = document.getElementById('selectedCurrency');
    if (currencySpan) {
        currencySpan.textContent = selectedCurrency;
    }
}

function addAddress() {
    if (addressCount >= 7) {
        alert('Maximum 7 addresses allowed');
        return;
    }
    
    addressCount++;
    const container = document.querySelector('.addresses-container');
    const newAddressGroup = document.createElement('div');
    newAddressGroup.className = 'address-input-group';
    
    // Calculate default percentage
    const defaultPercentage = Math.floor(100 / addressCount);
    
    newAddressGroup.innerHTML = `
        <label>Address ${addressCount}</label>
        <div class="address-input">
            <input type="text" placeholder="Enter wallet address" required>
            <input type="number" placeholder="%" min="1" max="100" value="${defaultPercentage}" class="percentage">
        </div>
    `;
    
    container.appendChild(newAddressGroup);
    
    // Redistribute percentages
    redistributePercentages();
    
    // Hide add button if we've reached the limit
    if (addressCount >= 7) {
        document.querySelector('.btn-add-address').style.display = 'none';
    }
}

function redistributePercentages() {
    const percentageInputs = document.querySelectorAll('.percentage');
    const equalPercentage = Math.floor(100 / percentageInputs.length);
    
    percentageInputs.forEach((input, index) => {
        if (index === percentageInputs.length - 1) {
            // Last input gets the remainder
            const remainder = 100 - (equalPercentage * (percentageInputs.length - 1));
            input.value = remainder;
        } else {
            input.value = equalPercentage;
        }
    });
}

function updatePercentages() {
    const percentageInputs = document.querySelectorAll('.percentage');
    let total = 0;
    
    percentageInputs.forEach(input => {
        total += parseFloat(input.value) || 0;
    });
    
    // Visual feedback if total doesn't equal 100%
    percentageInputs.forEach(input => {
        if (total !== 100) {
            input.style.borderColor = '#ef4444';
        } else {
            input.style.borderColor = '';
        }
    });
}

function validateAddresses() {
    const addressInputs = document.querySelectorAll('.address-input input[type="text"]');
    const nextBtn = document.querySelector('[data-step="3"] .btn-next');
    
    let allValid = true;
    addressInputs.forEach(input => {
        if (!input.value.trim()) {
            allValid = false;
        }
    });
    
    nextBtn.disabled = !allValid;
}

function nextStep() {
    if (currentStep < 5) {
        // Hide current step
        document.querySelector(`[data-step="${currentStep}"]`).classList.remove('active');
        
        // Show next step
        currentStep++;
        document.querySelector(`[data-step="${currentStep}"]`).classList.add('active');
        
        // Update order summary if moving to step 4
        if (currentStep === 4) {
            updateOrderSummary();
            generateQRCode();
        }
    }
}

function prevStep() {
    if (currentStep > 1) {
        // Hide current step
        document.querySelector(`[data-step="${currentStep}"]`).classList.remove('active');
        
        // Show previous step
        currentStep--;
        document.querySelector(`[data-step="${currentStep}"]`).classList.add('active');
    }
}

function updateOrderSummary() {
    const btcAmount = parseFloat(document.getElementById('btcAmount').value) || 0;
    const usdAmount = btcAmount * btcPrice;
    const serviceFee = usdAmount * 0.01; // 1% fee
    const receiveAmount = btcAmount * cryptoPrices[selectedCurrency];
    
    document.getElementById('summaryBTC').textContent = `${btcAmount} BTC`;
    document.getElementById('summaryReceive').textContent = `~${receiveAmount.toFixed(6)} ${selectedCurrency}`;
    document.getElementById('summaryFee').textContent = `~$${serviceFee.toFixed(2)}`;
}

function generatePaymentAddress() {
    // Generate a realistic-looking Bitcoin address
    const addressTypes = [
        'bc1q', // Bech32
        '3', // P2SH
        '1' // P2PKH
    ];
    
    const randomType = addressTypes[Math.floor(Math.random() * addressTypes.length)];
    let address;
    
    if (randomType === 'bc1q') {
        // Bech32 address
        address = 'bc1q' + generateRandomString(39, 'abcdefghijklmnopqrstuvwxyz0123456789');
    } else if (randomType === '3') {
        // P2SH address
        address = '3' + generateRandomString(33, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz123456789');
    } else {
        // P2PKH address
        address = '1' + generateRandomString(33, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz123456789');
    }
    
    document.getElementById('paymentAddress').textContent = address;
}

function generateRandomString(length, charset) {
    let result = '';
    for (let i = 0; i < length; i++) {
        result += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return result;
}

function generateQRCode() {
    const address = document.getElementById('paymentAddress').textContent;
    const btcAmount = parseFloat(document.getElementById('btcAmount').value) || 0;
    
    // Bitcoin URI format
    const bitcoinURI = `bitcoin:${address}?amount=${btcAmount}`;
    
    // Clear previous QR code
    document.getElementById('qrcode').innerHTML = '';
    
    // Generate new QR code
    QRCode.toCanvas(document.getElementById('qrcode'), bitcoinURI, {
        width: 200,
        height: 200,
        colorDark: '#000000',
        colorLight: '#ffffff',
        margin: 2
    }, function(error) {
        if (error) {
            console.error('QR Code generation failed:', error);
            // Fallback: create a simple QR code with just the address
            QRCode.toCanvas(document.getElementById('qrcode'), address, {
                width: 200,
                height: 200,
                colorDark: '#000000',
                colorLight: '#ffffff',
                margin: 2
            });
        }
    });
}

function copyAddress() {
    const address = document.getElementById('paymentAddress').textContent;
    
    if (navigator.clipboard) {
        navigator.clipboard.writeText(address).then(function() {
            showCopyFeedback();
        }).catch(function(err) {
            console.error('Failed to copy: ', err);
            fallbackCopy(address);
        });
    } else {
        fallbackCopy(address);
    }
}

function fallbackCopy(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();
    
    try {
        document.execCommand('copy');
        showCopyFeedback();
    } catch (err) {
        console.error('Fallback copy failed:', err);
        alert('Copy failed. Please copy manually: ' + text);
    }
    
    document.body.removeChild(textArea);
}

function showCopyFeedback() {
    const copyBtn = document.querySelector('.copy-btn');
    const originalHTML = copyBtn.innerHTML;
    
    copyBtn.innerHTML = '<i class="fas fa-check"></i>';
    copyBtn.style.background = '#10b981';
    
    setTimeout(function() {
        copyBtn.innerHTML = originalHTML;
        copyBtn.style.background = '';
    }, 2000);
}

function startMixing() {
    // Simulate starting the mixing process
    nextStep(); // Move to step 5 (processing)
    
    // Generate a random transaction ID
    const txId = 'CMX-' + generateRandomString(16, 'abcdefghijklmnopqrstuvwxyz0123456789');
    document.getElementById('txId').textContent = txId;
    
    // Simulate processing steps
    setTimeout(function() {
        updateProgressStep(2); // Converting Currency
    }, 3000);
    
    setTimeout(function() {
        updateProgressStep(3); // Sending to Addresses
        showCompletionMessage();
    }, 8000);
}

function updateProgressStep(stepNumber) {
    const steps = document.querySelectorAll('.progress-step');
    
    // Remove active from all steps
    steps.forEach(step => step.classList.remove('active'));
    
    // Add completed to previous steps and active to current
    steps.forEach((step, index) => {
        if (index + 1 < stepNumber) {
            step.classList.add('completed');
        } else if (index + 1 === stepNumber) {
            step.classList.add('active');
        }
    });
}

function showCompletionMessage() {
    const processingStatus = document.querySelector('.processing-status');
    processingStatus.innerHTML = `
        <div class="processing-icon">
            <i class="fas fa-check-circle" style="color: #10b981;"></i>
        </div>
        <h3>Transaction Complete!</h3>
        <p>Your cryptocurrency has been successfully sent to the specified addresses. Thank you for using CryptoMix!</p>
        
        <div class="progress-steps">
            <div class="progress-step completed">
                <i class="fas fa-check"></i>
                <span>Bitcoin Received</span>
            </div>
            <div class="progress-step completed">
                <i class="fas fa-check"></i>
                <span>Currency Converted</span>
            </div>
            <div class="progress-step completed">
                <i class="fas fa-check"></i>
                <span>Sent to Addresses</span>
            </div>
        </div>

        <div class="transaction-id">
            <strong>Transaction ID:</strong>
            <span id="txId">${document.getElementById('txId').textContent}</span>
        </div>
        
        <button type="button" class="btn-next" onclick="startNewTransaction()" style="margin-top: 2rem;">
            Start New Transaction
        </button>
    `;
}

function startNewTransaction() {
    // Reset the form
    currentStep = 1;
    selectedCurrency = '';
    addressCount = 1;
    
    // Reset form values
    document.getElementById('btcAmount').value = '';
    document.getElementById('amountUSD').textContent = '≈ $0 USD';
    
    // Clear currency selection
    document.querySelectorAll('.currency-option').forEach(opt => {
        opt.classList.remove('selected');
    });
    
    // Reset addresses
    const addressesContainer = document.querySelector('.addresses-container');
    addressesContainer.innerHTML = `
        <div class="address-input-group">
            <label>Address 1</label>
            <div class="address-input">
                <input type="text" placeholder="Enter wallet address" required>
                <input type="number" placeholder="%" min="1" max="100" value="100" class="percentage">
            </div>
        </div>
    `;
    
    // Show add address button
    document.querySelector('.btn-add-address').style.display = 'flex';
    
    // Reset step visibility
    document.querySelectorAll('.form-step').forEach(step => {
        step.classList.remove('active');
    });
    document.querySelector('[data-step="1"]').classList.add('active');
    
    // Generate new payment address
    generatePaymentAddress();
    
    // Re-setup event listeners for new elements
    setupEventListeners();
}

// Smooth scrolling for navigation links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Add some random price fluctuation to make it more realistic
setInterval(function() {
    // Simulate small BTC price changes (±1%)
    const fluctuation = (Math.random() - 0.5) * 0.02; // ±1%
    btcPrice = btcPrice * (1 + fluctuation);
    
    // Update displayed amounts if on step 1
    if (currentStep === 1) {
        updateBTCAmount();
    }
}, 30000); // Update every 30 seconds

// Add loading animation to buttons
document.addEventListener('click', function(e) {
    if (e.target.matches('.btn-next, .btn-submit')) {
        if (!e.target.disabled) {
            const originalText = e.target.textContent;
            e.target.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Loading...';
            
            setTimeout(function() {
                e.target.textContent = originalText;
            }, 1000);
        }
    }
});

// Add floating animation to features
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver(function(entries) {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe feature cards
document.querySelectorAll('.feature').forEach(feature => {
    feature.style.opacity = '0';
    feature.style.transform = 'translateY(30px)';
    feature.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    observer.observe(feature);
});

// Add particle animation effect
function createParticle() {
    const particle = document.createElement('div');
    particle.style.position = 'fixed';
    particle.style.width = (2 * (0.5 + Math.random())).toFixed(1) + 'px';
    particle.style.height = (2 * (0.5 + Math.random())).toFixed(1) + 'px';
    particle.style.background = '#6366f1';
    particle.style.borderRadius = '50%';
    particle.style.pointerEvents = 'none';
    particle.style.zIndex = '1';
    particle.style.opacity = (0.7 * (0.5 + Math.random())).toFixed(2);
    
    // Random starting position
    particle.style.left = Math.random() * window.innerWidth + 'px';
    particle.style.top = window.innerHeight + 'px';
    
    document.body.appendChild(particle);
    
    // Animate upward
    const animation = particle.animate([
        { transform: 'translateY(0px)', opacity: 0.1 + Math.random() * 0.6 },
        { transform: `translateY(-${window.innerHeight + 100}px)`, opacity: 0 }
    ], {
        duration: 3000 + Math.random() * 2000,
        easing: 'linear'
    });
    
    animation.onfinish = () => {
        particle.remove();
    };
}

// Create particles periodically
setInterval(createParticle, 5);
