// static/js/notifications.js
const showNotification = (message, type = 'info') => {
  const toast = document.createElement('div');
  toast.className = `alert alert-${type} mb-2`;
  toast.innerHTML = `
    <span>${message}</span>
    <button onclick="this.parentElement.remove()" class="btn btn-ghost btn-xs">✕</button>
  `;
  
  const toastContainer = document.querySelector('.toast');
  toastContainer.appendChild(toast);
  
  setTimeout(() => toast.remove(), 5000);
};
