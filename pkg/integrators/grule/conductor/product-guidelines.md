# Product Guidelines

## 1. Visual Identity & Design System
The Grule Engine frontend follows a modern, enterprise-tech aesthetic characterized by a cohesive "Blue-Violet Gradient" theme. This design language emphasizes clarity, professionalism, and forensic-grade detail.

*   **Color Palette:**
    *   **Primary Accent:** A dominant blue-violet linear gradient (`#667eea` → `#764ba2`) is used extensively on headers, search panels, rule selectors, modal titles, and primary action buttons.
    *   **Backgrounds:** The main application body utilizes a light grayish tone (`#f7fafc`) or subtle gradients to reduce eye strain. Content areas are distinct white cards.
    *   **Containers:** White backgrounds with soft box shadows (`0 4px–8px rgba(0,0,0,0.1)`) and rounded corners (8–12px radius) create a clean, card-based layout.
*   **Typography:** The system uses a standard, high-readability sans-serif stack: `'Segoe UI', Tahoma, Geneva, Verdana, sans-serif`.
*   **Iconography:** A rich combination of **Bootstrap Icons** (`bi-*`) and **Font Awesome** (`fas fa-*`) provides visual cues throughout the interface.
*   **Components:**
    *   **Buttons:** Bootstrap 5 style with embedded icons, utilizing the primary gradient or standard success/danger/warning/info variants.
    *   **Modals:** Large (`modal-xl`), scrollable modals with dark gradient headers are used for detailed inspections.
    *   **Status Indicators:** Heavy use of Bootstrap badges (success, danger, warning, info, secondary) ensures quick visual status communication.

## 2. User Experience (UX) & Interaction
*   **Feedback & Error Handling:**
    *   **Technical Depth:** Detailed error logs and raw data are exposed in modals, catering to the technical audience (developers, sysadmins) who need deep forensic insight.
    *   **Status Clarity:** Color-coded badges and clear visual indicators are used to communicate system state instantly.
*   **Tone & Voice:**
    *   **Professional & Precise:** The interface text and documentation maintain a professional, technical, and precise tone, suitable for operational security tools.
*   **Accessibility:**
    *   **Standards:** The design aims for adherence to **WCAG 2.1 AA** standards to ensure inclusivity.

## 3. Data Visualization
Given the system's focus on forensics and auditability, data visualization is a critical component:

*   **Tabular Data:** **jqGrid** and **DataTables** are used for high-density, interactive data presentation, often featuring headers styled with the signature blue-violet gradient and row hover effects.
*   **Sequential Analysis ("Movie Mode"):** Timelines are used to visualize the sequence of events, allowing users to replay rule executions frame-by-frame.
*   **Trend Analysis:** Charts and graphs are employed where appropriate to show trends and aggregate metrics.
