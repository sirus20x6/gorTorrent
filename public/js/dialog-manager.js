// /public/js/dialog-manager.js

class DialogManager {
    constructor() {
        this.maxZ = 2000;
        this.visible = [];
        this.items = {};
        this.modalState = false;
        
        // Close on escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && !e.target.matches('input, textarea')) {
                this.hideTopmost();
            }
        });

        // Setup modal backdrop if it doesn't exist
        if (!document.getElementById('modalbg')) {
            const backdrop = document.createElement('div');
            backdrop.id = 'modalbg';
            backdrop.className = 'fixed inset-0 bg-black/50 hidden';
            backdrop.style.zIndex = this.maxZ - 1;
            document.body.appendChild(backdrop);
        }
    }

    make(id, name, content, isModal = false, noClose = false) {
        // Create dialog container if it doesn't exist
        if (!document.getElementById('dialog-container')) {
            const container = document.createElement('div');
            container.id = 'dialog-container';
            document.body.appendChild(container);
        }

        // Convert content if it's a string
        if (typeof content === 'string') {
            content = new DOMParser().parseFromString(content, 'text/html').body.firstChild;
        }

        // Create dialog HTML
        const dialog = document.createElement('div');
        dialog.id = id;
        dialog.className = 'dlg-window';
        dialog.style.display = 'none';

        const header = document.createElement('div');
        header.className = 'dlg-header';
        
        const headerText = document.createElement('div');
        headerText.id = `${id}-header`;
        headerText.textContent = name;
        
        const closeBtn = document.createElement('a');
        closeBtn.href = '#';
        closeBtn.className = 'dlg-close';

        header.appendChild(headerText);
        if (!noClose) {
            header.appendChild(closeBtn);
        }

        dialog.appendChild(header);
        dialog.appendChild(content);

        document.getElementById('dialog-container').appendChild(dialog);
        return this.add(id, isModal, noClose);
    }

    add(id, isModal = false, noClose = false) {
        const dialog = document.getElementById(id);
        if (!dialog) return this;

        dialog.dataset.modal = isModal;
        dialog.dataset.nokeyclose = noClose;

        // Setup close button if not disabled
        if (!noClose) {
            const closeBtn = dialog.querySelector('.dlg-close');
            if (closeBtn) {
                closeBtn.addEventListener('click', (e) => {
                    e.preventDefault();
                    this.hide(id);
                });
            }
        }

        // Setup button handlers
        dialog.querySelectorAll('.Cancel').forEach(btn => {
            btn.addEventListener('click', () => this.hide(id));
        });

        dialog.querySelectorAll('.Button').forEach(btn => {
            btn.addEventListener('focus', function() { this.blur(); });
        });

        // Setup dialog focus and keyboard handling
        dialog.tabIndex = 0;
        dialog.addEventListener('mousedown', (e) => {
            if (!this.modalState) {
                this.bringToTop(id);
            }
        });

        dialog.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && 
                !e.target.matches('textarea') && 
                !dialog.querySelector('.OK')?.disabled) {
                dialog.querySelector('.OK')?.click();
            }
        });

        // Initialize drag functionality
        this.setupDrag(id);

        // Store dialog handlers
        this.items[id] = {
            beforeShow: null,
            afterShow: null,
            beforeHide: null,
            afterHide: null
        };

        return this;
    }

    setupDrag(id) {
        const dialog = document.getElementById(id);
        if (!dialog) return;

        const header = dialog.querySelector('.dlg-header');
        if (!header) return;

        let isDragging = false;
        let currentX, currentY, initialX, initialY, xOffset = 0, yOffset = 0;
        
        const dragStart = (e) => {
            if (window.innerWidth < 768) return; // Disable drag on mobile
            if (e.target.matches('a, button')) return; // Don't drag from buttons
            
            initialX = e.clientX - xOffset;
            initialY = e.clientY - yOffset;
            
            if (e.target === header || e.target.parentElement === header) {
                isDragging = true;
                document.body.style.cursor = 'grabbing';
            }
        };

        const drag = (e) => {
            if (!isDragging) return;
            e.preventDefault();

            currentX = e.clientX - initialX;
            currentY = e.clientY - initialY;

            // Constrain to window bounds
            const bounds = dialog.getBoundingClientRect();
            const maxX = window.innerWidth - bounds.width;
            const maxY = window.innerHeight - bounds.height;

            currentX = Math.min(Math.max(0, currentX), maxX);
            currentY = Math.min(Math.max(0, currentY), maxY);

            xOffset = currentX;
            yOffset = currentY;

            dialog.style.transform = `translate(${currentX}px, ${currentY}px)`;
        };

        const dragEnd = () => {
            isDragging = false;
            document.body.style.cursor = '';
        };

        header.addEventListener('mousedown', dragStart);
        document.addEventListener('mousemove', drag);
        document.addEventListener('mouseup', dragEnd);
    }

    center(id) {
        const dialog = document.getElementById(id);
        if (!dialog) return;

        const bounds = dialog.getBoundingClientRect();
        const x = Math.max((window.innerWidth - bounds.width) / 2, 0);
        const y = Math.max((window.innerHeight - bounds.height) / 2, 0);
        
        dialog.style.transform = `translate(${x}px, ${y}px)`;
    }

    toggle(id) {
        const pos = this.visible.indexOf(id);
        if (pos >= 0) {
            this.hide(id);
        } else {
            this.show(id);
        }
    }

    setModalState() {
        const backdrop = document.getElementById('modalbg');
        if (backdrop) {
            backdrop.style.display = 'block';
            this.bringToTop('modalbg');
        }
        this.modalState = true;
    }

    clearModalState() {
        const backdrop = document.getElementById('modalbg');
        if (backdrop) {
            backdrop.style.display = 'none';
        }
        this.modalState = false;
    }

    show(id, callback) {
        const dialog = document.getElementById(id);
        if (!dialog) return;

        // Handle mobile specific behaviors
        if (window.innerWidth < 768) {
            // Close any bootstrap components that might interfere
            const sidepanel = document.getElementById('offcanvas-sidepanel');
            const topMenu = document.getElementById('top-menu');
            if (sidepanel && window.bootstrap?.Offcanvas) {
                const instance = bootstrap.Offcanvas.getInstance(sidepanel);
                instance?.hide();
            }
            if (topMenu && window.bootstrap?.Collapse) {
                const instance = bootstrap.Collapse.getInstance(topMenu);
                instance?.hide();
            }
        }

        if (dialog.dataset.modal === 'true') {
            this.setModalState();
        }

        // Execute before show handler
        if (this.items[id]?.beforeShow) {
            this.items[id].beforeShow(id);
        }

        this.center(id);
        dialog.style.display = 'block';

        // Execute after show handler
        if (this.items[id]?.afterShow) {
            this.items[id].afterShow(id);
        }

        this.bringToTop(id);
        
        if (callback) callback();
    }

    hide(id, callback) {
        const pos = this.visible.indexOf(id);
        if (pos >= 0) {
            this.visible.splice(pos, 1);
        }

        const dialog = document.getElementById(id);
        if (!dialog) return;

        // Execute before hide handler
        if (this.items[id]?.beforeHide) {
            this.items[id].beforeHide(id);
        }

        dialog.style.display = 'none';

        // Execute after hide handler
        if (this.items[id]?.afterHide) {
            this.items[id].afterHide(id);
        }

        if (dialog.dataset.modal === 'true') {
            this.clearModalState();
        }

        if (callback) callback();
    }

    setHandler(id, type, handler) {
        if (this.items[id]) {
            this.items[id][type] = handler;
        }
        return this;
    }

    addHandler(id, type, handler) {
        if (this.items[id]) {
            const existing = this.items[id][type];
            if (existing) {
                this.items[id][type] = function() {
                    existing();
                    handler();
                };
            } else {
                this.items[id][type] = handler;
            }
        }
        return this;
    }

    isModalState() {
        return this.modalState;
    }

    bringToTop(id) {
        if (this.items[id]) {
            const pos = this.visible.indexOf(id);
            if (pos >= 0) {
                if (pos === this.visible.length - 1) return;
                this.visible.splice(pos, 1);
            }
            this.visible.push(id);
        }

        const dialog = document.getElementById(id);
        if (dialog) {
            dialog.style.zIndex = ++this.maxZ;
            if (!navigator.userAgent.includes('Opera')) {
                dialog.focus();
            }
        }
    }

    hideTopmost() {
        if (this.visible.length) {
            const topmostId = this.visible[this.visible.length - 1];
            const dialog = document.getElementById(topmostId);
            if (dialog && dialog.dataset.nokeyclose !== 'true') {
                this.hide(topmostId);
                return true;
            }
        }
        return false;
    }
}

// Initialize singleton instance
const theDialogManager = new DialogManager();