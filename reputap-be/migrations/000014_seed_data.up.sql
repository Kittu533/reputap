-- Seed roles
INSERT INTO roles (name, description) VALUES
('super_admin', 'Full system access'),
('admin', 'Operational management access'),
('operator', 'Limited operational access');

-- Seed permissions
INSERT INTO permissions (name, description, module) VALUES
-- Auth
('auth.manage_users', 'Create and manage user accounts', 'auth'),
('auth.manage_roles', 'Manage roles and permissions', 'auth'),
-- Customers
('customers.view', 'View customers', 'customers'),
('customers.create', 'Create customers', 'customers'),
('customers.update', 'Update customers', 'customers'),
('customers.delete', 'Delete customers', 'customers'),
-- Stores
('stores.view', 'View stores', 'stores'),
('stores.create', 'Create stores', 'stores'),
('stores.update', 'Update stores', 'stores'),
('stores.manage_status', 'Change store status', 'stores'),
('stores.manage_url', 'Change review URL', 'stores'),
-- Orders
('orders.view', 'View orders', 'orders'),
('orders.create', 'Create orders', 'orders'),
('orders.update', 'Update orders', 'orders'),
('orders.manage_status', 'Change order status', 'orders'),
-- Products
('products.view', 'View products', 'products'),
('products.create', 'Create products', 'products'),
('products.update', 'Update products', 'products'),
-- QR Codes
('qrcodes.view', 'View QR codes', 'qrcodes'),
('qrcodes.generate', 'Generate QR codes', 'qrcodes'),
('qrcodes.manage_status', 'Change QR status', 'qrcodes'),
-- Production
('production.view', 'View production jobs', 'production'),
('production.create', 'Create production jobs', 'production'),
('production.update', 'Update production status', 'production'),
-- QC
('qc.view', 'View QC records', 'qc'),
('qc.create', 'Submit QC results', 'qc'),
-- Shipping
('shipping.view', 'View shipments', 'shipping'),
('shipping.create', 'Create shipments', 'shipping'),
('shipping.update', 'Update shipping status', 'shipping'),
-- Analytics
('analytics.view', 'View analytics', 'analytics'),
-- Audit
('audit.view', 'View audit logs', 'audit'),
-- Settings
('settings.manage', 'Manage system settings', 'settings');

-- Assign all permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- Assign most permissions to admin (exclude auth.manage_users, auth.manage_roles, settings.manage, audit.view)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE name NOT IN ('auth.manage_users', 'auth.manage_roles', 'settings.manage');

-- Assign limited permissions to operator
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE name IN (
    'orders.view',
    'production.view', 'production.update',
    'qc.view', 'qc.create',
    'shipping.view', 'shipping.update',
    'qrcodes.view',
    'stores.view',
    'customers.view',
    'products.view'
);

-- Seed default super admin (password: admin123)
-- bcrypt hash for 'admin123'
INSERT INTO users (name, email, password_hash, role_id, status) VALUES
('Super Admin', 'admin@reputap.id', '$2a$10$BWefn4Lz1nFfOA92nVSqFuOhUPCI/RXFxTEBz0l6G6hOxHPEBli7e', 1, 'active');
