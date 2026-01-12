const auditApp = Vue.createApp({
    data() {
        const now = new Date();
        const yesterday = new Date(now.getTime() - 24*60*60*1000);
        
        return {
            searchImei: '',
            searchStartDate: yesterday.toISOString().split('T')[0],
            searchEndDate: now.toISOString().split('T')[0],
            executions: [],
            selectedExecution: null,
            expandedSteps: {},
            searched: false,
            loading: false,
            error: null,
            recentImeis: [],
            recentImeiFiltered: []
        }
    },
    mounted() {
        // Cargar IMEIs recientes al montar
        this.loadRecentImeis();
    },
    methods: {
        formatDate(dateString) {
            const date = new Date(dateString);
            return date.toLocaleString();
        },
        formatValue(val) {
            if (typeof val === 'object') {
                return JSON.stringify(val).substring(0, 50) + '...';
            }
            return String(val).substring(0, 100);
        },
        loadRecentImeis() {
            fetch('/api/audit/imeis/recent?limit=20')
                .then(r => {
                    if (!r.ok) throw new Error(`HTTP ${r.status}`);
                    return r.json();
                })
                .then(data => {
                    this.recentImeis = data.imeis || [];
                })
                .catch(err => {
                    console.error('Error loading recent IMEIs:', err);
                });
        },
        filterRecentImeis() {
            if (this.searchImei.trim() === '') {
                this.recentImeiFiltered = [];
                return;
            }
            
            const query = this.searchImei.toLowerCase();
            this.recentImeiFiltered = this.recentImeis.filter(item => 
                item.imei.toLowerCase().includes(query)
            );
        },
        selectImei(imei) {
            this.searchImei = imei;
            this.recentImeiFiltered = [];
            this.searchExecutions();
        },
        searchExecutions() {
            if (!this.searchImei.trim()) {
                alert('Please enter an IMEI');
                return;
            }

            this.loading = true;
            this.error = null;
            this.searched = true;

            const params = new URLSearchParams({
                imei: this.searchImei,
                start_date: this.searchStartDate,
                end_date: this.searchEndDate
            });

            fetch(`/api/audit/executions?${params}`)
                .then(r => {
                    if (!r.ok) throw new Error(`HTTP ${r.status}`);
                    return r.json();
                })
                .then(data => {
                    this.executions = data.executions || [];
                    this.loading = false;
                    if (this.executions.length === 0) {
                        this.error = 'No executions found for the selected criteria';
                    }
                })
                .catch(err => {
                    console.error('Error searching:', err);
                    this.error = `Error: ${err.message}`;
                    this.loading = false;
                });
        },
        viewExecution(executionId) {
            this.loading = true;
            fetch(`/api/audit/executions/${executionId}`)
                .then(r => {
                    if (!r.ok) throw new Error(`HTTP ${r.status}`);
                    return r.json();
                })
                .then(data => {
                    this.selectedExecution = data.trace || data;
                    this.expandedSteps = {}; // Reset expanded steps
                    this.loading = false;
                })
                .catch(err => {
                    console.error('Error fetching detail:', err);
                    this.error = `Error: ${err.message}`;
                    this.loading = false;
                });
        },
        toggleStepDetail(idx) {
            // Vue 3 way: just toggle the property
            if (this.expandedSteps[idx]) {
                delete this.expandedSteps[idx];
            } else {
                this.$set(this.expandedSteps, idx, true);
            }
        }
    }
});

auditApp.mount('#app');
