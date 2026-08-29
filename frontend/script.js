let activeOutputFile = "";

// Toggle specific inputs depending on chosen task type
function renderOperationUI() {
    const type = document.getElementById("operationType").value;
    const sections = document.querySelectorAll(".dynamic-section");
    sections.forEach(sec => sec.style.display = "none");

    const submitBtn = document.getElementById("submitBtn");

    if (type) {
        const targetSection = document.getElementById(`opts-${type}`);
        if (targetSection) targetSection.style.display = "block";
        submitBtn.disabled = false;
    } else {
        submitBtn.disabled = true;
    }
}

// Process file upload first (if present), then trigger POST /task
async function submitTask() {
    const type = document.getElementById("operationType").value;
    let payloadData = {};

    let CFPayload = {};
    let fileInput = null;

    if (type === "compress") {
        fileInput = document.getElementById("file-compress");
        payloadData.level = document.getElementById("ratio").value;
    } else if (type === "watermark") {
        fileInput = document.getElementById("file-watermark");
        payloadData.text = document.getElementById("watermarkText").value;
    } else if (type === "convert_format") {
        fileInput = document.getElementById("file-convert");
        payloadData.target_format = document.getElementById("targetFormat").value;
        CFPayload.output_format = document.getElementById("targetFormat").value;
    }

    // Step 1: Upload file if selected
    if (fileInput && fileInput.files.length > 0) {
        const file = fileInput.files[0];
        const formData = new FormData();
        formData.append("newFile", file);

        try {
            const uploadRes = await fetch("/upload", {
                method: "POST",
                body: formData
            });
            if (!uploadRes.ok) throw new Error("File upload failed");
            payloadData.filename = file.name;
            CFPayload.file_name = file.name;
        } catch (err) {
            alert("Error uploading file: " + err.message);
            return;
        }
    }

    // Step 2: Post Task metadata
    try {
        const taskRes = await fetch("/task", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                type: type,
                payload: CFPayload
            })
        });

        if (!taskRes.ok) throw new Error("Failed to create task");

        const data = await taskRes.json();
        document.getElementById("trackTaskId").value = data.task_id;
        alert(`Task submitted successfully! ID: ${data.task_id}`);
        checkStatus(); // Auto-trigger status check
    } catch (err) {
        alert("Error creating task: " + err.message);
    }
}

// Query GET /task/{task_id}
async function checkStatus() {
    const taskId = document.getElementById("trackTaskId").value.trim();
    if (!taskId) {
        alert("Please enter a Task ID");
        return;
    }

    try {
        const res = await fetch(`/task/${taskId}`);
        if (!res.ok) {
            if (res.status === 404) alert("Task not found");
            return;
        }

        const data = await res.json();

        // Populate UI
        const statusResult = document.getElementById("statusResult");
        const statusBadge = document.getElementById("statusBadge");
        const createdAt = document.getElementById("createdAt");
        const downloadBtn = document.getElementById("downloadBtn");
        const errorContainer = document.getElementById("errorContainer");

        statusResult.classList.remove("hidden");
        statusBadge.textContent = data.status;
        statusBadge.className = `status-badge status-${data.status}`;
        createdAt.textContent = data.created_at || "N/A";

        if (data.status === "completed") {
            activeOutputFile = data.output_file;
            downloadBtn.classList.remove("hidden");
            errorContainer.classList.add("hidden");
        } else if (data.status === "failed") {
            document.getElementById("errorMsg").textContent = data.error || "Unknown error";
            errorContainer.classList.remove("hidden");
            downloadBtn.classList.add("hidden");
        } else {
            downloadBtn.classList.add("hidden");
            errorContainer.classList.add("hidden");
        }
    } catch (err) {
        alert("Error fetching task status: " + err.message);
    }
}

// Trigger file download via GET /download/{filename}
function downloadFile() {
    if (activeOutputFile) {
        window.location.href = `/download/${encodeURIComponent(activeOutputFile)}`;
    }
}