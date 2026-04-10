import React, { useEffect, useState } from 'react';
import { HardDrive, Bell, LogIn, ExternalLink, Filter } from 'lucide-react';

interface Product {
  id: string;
  asin: string;
  title: string;
  capacity_gb?: number;
  ram_capacity_gb?: number;
  form_factor?: string;
  ram_form_factor?: string;
  storage_interface?: string;
  ram_type?: string;
  price?: number;
  currency?: string;
  price_per_gb?: number;
  amazon_url?: string;
}

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  // Filters
  const [typeFilter, setTypeFilter] = useState<'all' | 'storage' | 'ram'>('all');

  useEffect(() => {
    fetch('/api/products')
      .then(res => res.json())
      .then(data => {
        if (data.products) setProducts(data.products);
        setLoading(false);
      })
      .catch(err => {
        console.error("Failed to fetch products", err);
        setLoading(false);
      });
  }, []);

  const formatPrice = (price?: number, currency?: string) => {
    if (price === undefined) return 'N/A';
    return new Intl.NumberFormat('en-GB', { style: 'currency', currency: currency || 'GBP' }).format(price);
  };

  const getCapacity = (p: Product) => p.capacity_gb || p.ram_capacity_gb || 0;
  
  const getCapacityLabel = (gb: number) => {
    if (gb >= 1000) return `${(gb / 1000).toFixed(1)} TB`;
    return `${gb} GB`;
  };

  const getTypeStr = (p: Product) => {
    if (p.ram_type) return `RAM - ${p.ram_type} ${p.ram_form_factor || ''}`;
    return `Storage - ${p.form_factor || ''} ${p.storage_interface || ''}`;
  };

  const filteredProducts = products.filter(p => {
    if (typeFilter === 'storage' && !p.capacity_gb) return false;
    if (typeFilter === 'ram' && !p.ram_capacity_gb) return false;
    return true;
  });

  return (
    <div className="app-container">
      {/* Navbar */}
      <nav className="navbar">
        <div className="logo">
          <HardDrive size={24} color="#3b82f6" />
          Disky
        </div>
        <div className="nav-actions">
          <button className="btn">
            <Bell size={16} /> Alerts
          </button>
          <button className="btn btn-primary">
            <LogIn size={16} /> Login
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className="main-content">
        {/* Sidebar */}
        <aside className="sidebar">
          <div className="filter-group">
            <h3><Filter size={14} style={{ display: 'inline', marginRight: 4 }} /> Filters</h3>
            <label className="filter-label">
              <input 
                type="radio" name="type" 
                checked={typeFilter === 'all'} 
                onChange={() => setTypeFilter('all')} 
              /> All
            </label>
            <label className="filter-label">
              <input 
                type="radio" name="type" 
                checked={typeFilter === 'storage'} 
                onChange={() => setTypeFilter('storage')} 
              /> Storage (SSD/HDD)
            </label>
            <label className="filter-label">
              <input 
                type="radio" name="type" 
                checked={typeFilter === 'ram'} 
                onChange={() => setTypeFilter('ram')} 
              /> Memory (RAM)
            </label>
          </div>
          
          <div className="filter-group" style={{ marginTop: '1rem' }}>
             <h3>More filters coming soon</h3>
             <p style={{ color: 'var(--text-muted)', fontSize: '0.85rem' }}>
                Advanced capacities, form factors, and speeds...
             </p>
          </div>
        </aside>

        {/* Data Table */}
        <section className="data-view">
          <div className="data-table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Capacity</th>
                  <th>Type</th>
                  <th>Price / GB</th>
                  <th>Total Price</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr><td colSpan={5} style={{ textAlign: 'center' }}>Loading impressive deals...</td></tr>
                ) : filteredProducts.map(p => (
                  <tr key={p.id}>
                    <td>
                      <a href={p.amazon_url} target="_blank" rel="noreferrer" className="product-link">
                        {p.title.length > 60 ? p.title.substring(0, 60) + '...' : p.title}
                        <ExternalLink size={14} style={{ opacity: 0.5 }} />
                      </a>
                    </td>
                    <td>
                      <span className="badge badge-info">{getCapacityLabel(getCapacity(p))}</span>
                    </td>
                    <td style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>
                      {getTypeStr(p)}
                    </td>
                    <td>
                      {p.price_per_gb ? (
                         <span className="badge badge-success">{formatPrice(p.price_per_gb, p.currency)}/GB</span>
                      ) : '-'}
                    </td>
                    <td className="price">
                      {formatPrice(p.price, p.currency)}
                    </td>
                  </tr>
                ))}
                {!loading && filteredProducts.length === 0 && (
                  <tr><td colSpan={5} style={{ textAlign: 'center' }}>No products found.</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;
