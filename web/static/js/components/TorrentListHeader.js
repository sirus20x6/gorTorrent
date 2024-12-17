// static/js/components/TorrentListHeader.js

const TorrentListHeader = () => {
    const [filter, setFilter] = React.useState('all');
    const [search, setSearch] = React.useState('');

    const filters = [
        { id: 'all', label: 'All' },
        { id: 'downloading', label: 'Downloading' },
        { id: 'seeding', label: 'Seeding' },
        { id: 'completed', label: 'Completed' },
        { id: 'active', label: 'Active' },
        { id: 'inactive', label: 'Inactive' },
        { id: 'error', label: 'Error' }
    ];

    const handleFilterChange = (newFilter) => {
        setFilter(newFilter);
        // Send HTMX request to update list
        htmx.ajax('GET', `/torrents?filter=${newFilter}`, {
            target: '#torrent-list',
            swap: 'innerHTML'
        });
    };

    const handleSearch = (e) => {
        setSearch(e.target.value);
        // Debounce search requests
        clearTimeout(window.searchTimeout);
        window.searchTimeout = setTimeout(() => {
            htmx.ajax('GET', `/torrents?search=${e.target.value}`, {
                target: '#torrent-list',
                swap: 'innerHTML'
            });
        }, 300);
    };

    const handleSort = (field) => {
        htmx.ajax('GET', `/torrents?sort=${field}`, {
            target: '#torrent-list',
            swap: 'innerHTML'
        });
    };

    return (
        <div className="flex flex-col gap-4 p-4">
            {/* Search and Add Button Row */}
            <div className="flex flex-wrap gap-4 items-center">
                <div className="join flex-1">
                    <div className="join-item flex items-center bg-base-200 px-4">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </div>
                    <input 
                        type="text" 
                        placeholder="Search torrents..." 
                        className="input input-bordered join-item flex-1"
                        value={search}
                        onChange={handleSearch}
                    />
                </div>
                <button 
                    className="btn btn-primary"
                    onClick={() => document.getElementById('add_torrent_modal').showModal()}
                >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4" />
                    </svg>
                    Add Torrent
                </button>
            </div>

            {/* Filters Row */}
            <div className="flex flex-wrap gap-4">
                <div className="join">
                    {filters.map(f => (
                        <button
                            key={f.id}
                            className={`btn btn-sm join-item ${filter === f.id ? 'btn-active' : ''}`}
                            onClick={() => handleFilterChange(f.id)}
                        >
                            {f.label}
                        </button>
                    ))}
                </div>

                <div className="join ml-auto">
                    <div className="dropdown dropdown-end join-item">
                        <button tabIndex={0} className="btn btn-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
                            </svg>
                            Sort
                        </button>
                        <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-base-200 rounded-box w-52">
                            <li><button onClick={() => handleSort('name')}>Name</button></li>
                            <li><button onClick={() => handleSort('size')}>Size</button></li>
                            <li><button onClick={() => handleSort('progress')}>Progress</button></li>
                            <li><button onClick={() => handleSort('speed')}>Speed</button></li>
                            <li><button onClick={() => handleSort('added')}>Date Added</button></li>
                        </ul>
                    </div>
                </div>
            </div>
            
            {/* Stats Row */}
            <div className="stats shadow">
                <div className="stat">
                    <div className="stat-title">Download Speed</div>
                    <div className="stat-value text-success text-2xl">0 KB/s</div>
                </div>
                
                <div className="stat">
                    <div className="stat-title">Upload Speed</div>
                    <div className="stat-value text-info text-2xl">0 KB/s</div>
                </div>
                
                <div className="stat">
                    <div className="stat-title">Active Torrents</div>
                    <div className="stat-value text-2xl">0</div>
                </div>
            </div>
        </div>
    );
};

export default TorrentListHeader;