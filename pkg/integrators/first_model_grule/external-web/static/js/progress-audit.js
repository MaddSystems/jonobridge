$(document).ready(function() {
    const API_BASE_URL = "/grule/api/audit/progress";

    // 1. Load rules into selector
    $.ajax({
        url: `${API_BASE_URL}/rules`,
        type: 'GET',
        success: function(response) {
            if (response.success) {
                const selector = $('#rule-selector');
                selector.empty().append('<option value="">Seleccione una regla</option>');
                response.rules.forEach(rule => {
                    selector.append(`<option value="${rule.rule_name}">${rule.rule_name} (${rule.total_frames} frames, ${rule.total_imeis} IMEIs)</option>`);
                });
            }
        },
        error: function() {
            alert('Error cargando las reglas.');
        }
    });

    // 2. Initialize Level 1 Grid (IMEIs)
    $("#imeis-grid").jqGrid({
        datatype: "local", // Initially local, will be changed to JSON on button click
        colNames: ['IMEI', 'Max Step', 'Frames', 'Última Frame'],
        colModel: [
            { name: 'imei', index: 'imei', width: 150, sorttype: 'text' },
            { name: 'max_step', index: 'max_step', width: 90, sorttype: 'int', align: 'center' },
            { name: 'total_frames', index: 'total_frames', width: 90, sorttype: 'int', align: 'center' },
            { name: 'last_frame_time', index: 'last_frame_time', width: 150, sorttype: 'date', formatter: 'date', formatoptions: { newformat: 'Y-m-d H:i:s' } }
        ],
        rowNum: 15,
        rowList: [15, 30, 50],
        pager: '#imeis-pager',
        sortname: 'last_frame_time',
        viewrecords: true,
        sortorder: "desc",
        autowidth: true,
        height: 'auto',
        caption: "Progreso por IMEI",
        jsonReader: {
            root: "data",
            page: "page",
            total: "total",
            records: "records",
            repeatitems: false,
            id: "imei"
        },
        onSelectRow: function(rowid) {
            if (rowid) {
                const ruleName = $("#rule-selector").val();
                $("#modal-imei").text(rowid);
                $("#timeline-grid").jqGrid('setGridParam', {
                    url: `${API_BASE_URL}/timeline?imei=${rowid}&rule_name=${ruleName}`,
                    datatype: 'json',
                    page: 1
                }).trigger("reloadGrid");
                $('#timelineModal').modal('show');
            }
        }
    });

    // 3. Initialize Level 2 Grid (Timeline)
    $("#timeline-grid").jqGrid({
        datatype: "local",
        colNames: ['Timestamp', 'Step', 'Stage', 'Snapshot'],
        colModel: [
            { name: 'execution_time', index: 'execution_time', width: 150, sorttype: 'date', formatter: 'date', formatoptions: { newformat: 'Y-m-d H:i:s' } },
            { name: 'step_number', index: 'step_number', width: 90, sorttype: 'int', align: 'center' },
            { name: 'stage_reached', index: 'stage_reached', width: 150 },
            { name: 'snapshot', width: 100, align: 'center', sortable: false, formatter: function(cellvalue, options, rowObject) {
                return '<button class="btn btn-primary btn-sm view-snapshot">Ver JSON</button>';
            }}
        ],
        rowNum: 10,
        rowList: [10, 20, 50],
        pager: '#timeline-pager',
        sortname: 'execution_time',
        viewrecords: true,
        sortorder: "asc",
        autowidth: true,
        height: 'auto',
        caption: "Detalle de Frames",
        jsonReader: {
            root: "frames",
            page: "page",
            total: "total",
            records: "records",
            repeatitems: false,
            id: "id"
        },
        gridComplete: function() {
            $(".view-snapshot").on("click", function() {
                const rowid = $(this).closest("tr").attr("id");
                const rowData = $("#timeline-grid").jqGrid('getRowData', rowid);
                const snapshotData = JSON.parse(rowData.snapshot); // Snapshot is stored as a string
                $("#json-content").text(JSON.stringify(snapshotData, null, 2));
                $('#snapshotModal').modal('show');
            });
        }
    });

    // 4. Button click handler
    $("#load-imeis-btn").on("click", function() {
        const ruleName = $("#rule-selector").val();
        if (ruleName) {
            $("#imeis-grid").jqGrid('setGridParam', {
                url: `${API_BASE_URL}/summary?rule_name=${ruleName}`,
                datatype: 'json',
                page: 1
            }).trigger("reloadGrid");
        } else {
            alert("Por favor, seleccione una regla.");
        }
    });
});
