import os
import requests
import json
from datetime import datetime
from flask import Flask, render_template, request, redirect, url_for, flash, jsonify, make_response

app = Flask(__name__)
app.secret_key = os.urandom(24)

# Get the API base URL from environment variable or use a default
API_BASE_URL = os.environ.get("API_BASE_URL", "https://jonobridge.madd.com.mx/grule")

def load_templates():
    """Loads rule templates from .grl files in rules_templates directory"""
    templates = {}
    rules_dir = os.path.join(os.path.dirname(__file__), 'rules_templates')
    
    # Check if rules_templates directory exists
    if not os.path.exists(rules_dir):
        print(f"Warning: Rules directory {rules_dir} does not exist")
        return templates
    
    print(f"Loading templates from: {rules_dir}")
    
    # Load each .grl file as a template
    for filename in os.listdir(rules_dir):
        if filename.endswith('.grl'):
            template_key = filename[:-4]  # Remove .grl extension
            filepath = os.path.join(rules_dir, filename)
            
            try:
                with open(filepath, 'r', encoding='utf-8') as f:
                    grl_content = f.read()
                
                # Extract template metadata from comments at the top
                lines = grl_content.split('\n')
                name = template_key.replace('_', ' ').title()
                category = "Custom Rules"
                description = f"Loaded from {filename}"
                
                # Try to extract metadata from comments in first 10 lines
                for line in lines[:10]:
                    line = line.strip()
                    if line.startswith('// name:'):
                        name = line.replace('// name:', '').strip()
                    elif line.startswith('// category:'):
                        category = line.replace('// category:', '').strip()
                    elif line.startswith('// description:'):
                        description = line.replace('// description:', '').strip()
                
                # Special handling for known templates
                if template_key == 'gpsgate_expressions':
                    name = "GPSgate Memory Buffer - Sistema Expresiones"
                    category = "Seguridad Buffer Circular" 
                    description = "SISTEMA EXPRESIONES GPSgate: Replica exacta con buffer circular de 10 posiciones. Implementa sistema 'Coincidir TODO' con 5 expresiones: geofences, offline, y métricas buffer (velocidad 90min, GSM últimos 5). Performance 500x mejorada vs GPSgate original."
                
                templates[template_key] = {
                    "name": name,
                    "category": category,
                    "description": description,
                    "grl": grl_content
                }
                
                print(f"✅ Loaded template: {template_key} ({name})")
                
            except Exception as e:
                print(f"❌ Error loading {filepath}: {e}")
                continue
    
    # Add fallback template if no files found
    if not templates:
        print("⚠️ No .grl files found, using fallback template")
        templates = {
            "fallback": {
                "name": "Fallback Template",
                "category": "System",
                "description": "Fallback template when no .grl files are found. Add .grl files to rules_templates directory.",
                "grl": '''rule FallbackExample "Fallback rule example" salience 100 {
    when
        true
    then
        actions.Log("⚠️ No .grl templates found - add files to rules_templates directory");
}'''
            }
        }
    else:
        print(f"📄 Total templates loaded: {len(templates)}")
    
    return templates


@app.route('/')
def index():
    """Main page, lists all the rules."""
    try:
        response = requests.get(f"{API_BASE_URL}/api/rules")
        response.raise_for_status()
        data = response.json()
        rules = data.get('rules', []) if data else []
    except requests.exceptions.RequestException as e:
        flash(f"Error fetching rules: {e}", "error")
        rules = []
    return render_template('index.html', rules=rules if rules else [])

@app.route('/rule/new', methods=['GET', 'POST'])
def new_rule():
    """Page to create a new rule."""
    templates = load_templates()
    if request.method == 'POST':
        rule_data = {
            "name": request.form['name'],
            "grl": request.form['grl'],
            "priority": int(request.form.get('priority', 100)),
            "active": 'active' in request.form
        }
        try:
            response = requests.post(f"{API_BASE_URL}/api/rules", json=rule_data)
            response.raise_for_status()
            data = response.json()
            
            # Validar que la API respondió con success=true
            if data.get('success'):
                flash('Rule created successfully!', 'success')
                return redirect(url_for('index'))
            else:
                # La API retornó success=false (error de sintaxis, etc)
                error_msg = data.get('error', 'Unknown error')
                flash(f"Error creating rule: {error_msg}", "error")
        except requests.exceptions.RequestException as e:
            flash(f"Error creating rule: {e}", "error")
    return render_template('form.html', rule=None, title="New Rule", templates_json=templates)

@app.route('/rule/edit/<int:rule_id>', methods=['GET', 'POST'])
def edit_rule(rule_id):
    """Page to edit an existing rule."""
    templates = load_templates()
    if request.method == 'POST':
        rule_data = {
            "name": request.form['name'],
            "grl": request.form['grl'],
            "priority": int(request.form.get('priority', 100)),
            "active": 'active' in request.form
        }
        try:
            response = requests.put(f"{API_BASE_URL}/api/rules/{rule_id}", json=rule_data)
            response.raise_for_status()
            flash('Rule updated successfully!', 'success')
            return redirect(url_for('index'))
        except requests.exceptions.RequestException as e:
            flash(f"Error updating rule: {e}", "error")

    try:
        response = requests.get(f"{API_BASE_URL}/api/rules/{rule_id}")
        response.raise_for_status()
        rule = response.json().get('rule', {})
    except requests.exceptions.RequestException as e:
        flash(f"Error fetching rule {rule_id}: {e}", "error")
        return redirect(url_for('index'))

    return render_template('form.html', rule=rule, title="Edit Rule", templates_json=templates)

@app.route('/rule/delete/<int:rule_id>', methods=['POST'])
def delete_rule(rule_id):
    """Endpoint to delete a rule."""
    try:
        response = requests.delete(f"{API_BASE_URL}/api/rules/{rule_id}")
        response.raise_for_status()
        flash('Rule deleted successfully!', 'success')
    except requests.exceptions.RequestException as e:
        flash(f"Error deleting rule: {e}", "error")
    return redirect(url_for('index'))

@app.route('/api/validate', methods=['POST'])
def validate_rule_syntax():
    """API endpoint for the form to AJAX-validate GRL syntax."""
    grl_content = request.json.get('grl', '')
    try:
        response = requests.post(f"{API_BASE_URL}/api/validate", json={"grl": grl_content})
        return jsonify(response.json()), response.status_code
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500
        
@app.route('/api/reload', methods=['POST'])
def force_reload():
    """Endpoint to trigger a hot-reload of the rules engine."""
    try:
        response = requests.post(f"{API_BASE_URL}/api/reload")
        response.raise_for_status()
        flash('Rules engine reload triggered successfully!', 'success')
    except requests.exceptions.RequestException as e:
        flash(f"Error triggering reload: {e}", "error")
    return redirect(url_for('index'))


# Audit Dashboard Routes
@app.route('/audit')
def audit_dashboard():
    """Render audit dashboard for rule execution visualization."""
    return render_template('audit_summary.html')


@app.route('/invalid-packets')
def invalid_packets_dashboard():
    """Render invalid packets dashboard to view IMEIs with invalid traces."""
    return render_template('invalid_packets.html')


@app.route('/api/audit/executions', methods=['GET'])
def get_audit_executions():
    """Proxy audit executions endpoint to Go backend."""
    try:
        imei = request.args.get('imei')
        start_date = request.args.get('start_date')
        end_date = request.args.get('end_date')
        limit = request.args.get('limit', '100')
        
        if not imei:
            return jsonify({"error": "Missing required parameter: imei"}), 400
        
        params = {'imei': imei, 'limit': limit}
        if start_date:
            params['start_date'] = start_date
        if end_date:
            params['end_date'] = end_date
        
        response = requests.get(f"{API_BASE_URL}/api/audit/executions", params=params)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/audit/executions/<int:execution_id>', methods=['GET'])
def get_audit_execution_detail(execution_id):
    """Proxy audit execution detail endpoint to Go backend."""
    try:
        response = requests.get(f"{API_BASE_URL}/api/audit/executions/{execution_id}")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/audit/imeis/recent', methods=['GET'])
def get_recent_imeis():
    """Proxy recent IMEIs endpoint to Go backend."""
    try:
        limit = request.args.get('limit', '50')
        response = requests.get(f"{API_BASE_URL}/api/audit/imeis/recent?limit={limit}")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/audit/imeis/search', methods=['GET'])
def search_imeis():
    """Proxy IMEI search endpoint to Go backend."""
    try:
        q = request.args.get('q', '')
        limit = request.args.get('limit', '50')
        response = requests.get(f"{API_BASE_URL}/api/audit/imeis/search?q={q}&limit={limit}")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/audit/grid', methods=['GET'])
def get_audit_grid():
    """Proxy audit grid endpoint to Go backend for jqGrid."""
    try:
        page = request.args.get('page', '1')
        rows = request.args.get('rows', '25')
        sidx = request.args.get('sidx', 'execution_date')
        sord = request.args.get('sord', 'DESC')
        searchText = request.args.get('searchText', '')
        
        params = {
            'page': page,
            'rows': rows,
            'sidx': sidx,
            'sord': sord,
            'searchText': searchText
        }
        
        response = requests.get(f"{API_BASE_URL}/api/audit/grid", params=params)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/invalid-packets', methods=['GET'])
def get_invalid_packets():
    """Proxy invalid packets endpoint to Go backend."""
    try:
        page = request.args.get('page', '1')
        rows = request.args.get('rows', '50')
        
        params = {
            'page': page,
            'rows': rows
        }
        
        response = requests.get(f"{API_BASE_URL}/api/invalid-packets", params=params)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/invalid-packets-csv', methods=['GET'])
def get_invalid_packets_csv():
    """Download invalid packets data as CSV."""
    try:
        # Obtener todos los datos sin paginación
        params = {
            'page': '1',
            'rows': '10000',  # Obtener muchos registros
            'searchText': request.args.get('searchText', ''),
            'sidx': request.args.get('sidx', 'last_seen')
        }
        
        response = requests.get(f"{API_BASE_URL}/api/invalid-packets", params=params)
        response.raise_for_status()
        data = response.json()
        
        # Crear CSV
        import io
        import csv
        
        output = io.StringIO()
        writer = csv.writer(output)
        
        # Escribir encabezados
        writer.writerow(['IMEI', 'Última Vez Visto', 'Ocurrencias'])
        
        # Escribir datos
        if 'rows' in data:
            for row in data['rows']:
                writer.writerow([
                    row.get('imei', ''),
                    row.get('last_seen', ''),
                    row.get('count', 0)
                ])
        
        # Preparar respuesta CSV
        output.seek(0)
        csv_data = output.getvalue()
        output.close()
        
        response = make_response(csv_data)
        response.headers['Content-Type'] = 'text/csv'
        response.headers['Content-Disposition'] = f'attachment; filename=invalid_packets_{datetime.now().strftime("%Y%m%d")}.csv'
        
        return response
        
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


# ========================== PROGRESS AUDIT ROUTES ==========================

@app.route('/progress-audit')
def progress_audit_dashboard():
    """Render progress audit dashboard."""
    return render_template('progress_audit.html')


@app.route('/progress-audit-movie')
def progress_audit_movie():
    """Render progress audit movie frames dashboard with jqGrid."""
    return render_template('progress_audit_movie.html')


@app.route('/api/progress/enable', methods=['POST'])
def enable_progress_audit():
    """Enable progress audit tracking."""
    try:
        response = requests.post(f"{API_BASE_URL}/api/audit/progress/enable")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/progress/disable', methods=['POST'])
def disable_progress_audit():
    """Disable progress audit tracking."""
    try:
        response = requests.post(f"{API_BASE_URL}/api/audit/progress/disable")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/progress/clear', methods=['POST'])
def clear_progress_audit():
    """Clear all progress audit data."""
    try:
        response = requests.post(f"{API_BASE_URL}/api/audit/progress/clear")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/progress/status', methods=['GET'])
def get_progress_status():
    """Get current progress audit status."""
    try:
        response = requests.get(f"{API_BASE_URL}/api/audit/progress/status")
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/progress/query', methods=['GET'])
def query_progress_audit():
    """Query progress audit data by IMEI."""
    try:
        imei = request.args.get('imei')
        limit = request.args.get('limit', '50')
        
        if not imei:
            return jsonify({"error": "IMEI parameter required"}), 400
        
        params = {'imei': imei, 'limit': limit}
        response = requests.get(f"{API_BASE_URL}/api/audit/progress", params=params)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


@app.route('/api/progress/timeline', methods=['GET'])
def get_progress_timeline():
    """Proxy progress timeline endpoint to Go backend."""
    try:
        imei = request.args.get('imei')
        rule_name = request.args.get('rule_name')
        limit = request.args.get('limit', '500')

        if not imei:
            return jsonify({"error": "IMEI parameter is required"}), 400

        params = {'imei': imei, 'limit': limit}
        if rule_name:
            params['rule_name'] = rule_name
            
        response = requests.get(f"{API_BASE_URL}/api/audit/progress/timeline", params=params)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5001)