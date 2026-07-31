var fileUploaded = false;
var fileName = "";
var outputFileName = "";

const fileInput = document.getElementById("fileInput");
const uploadBtn = document.getElementById("fileUploadButton");
const postTaskBtn = document.getElementById("postTaskButton");
const outputFileNameP = document.getElementById("outputFileName");

async function uploadFile() {
    if (fileInput.files.length === 0) {
        alert("Please select a file first");
        return;
    }

    const formData = new FormData();
    formData.append("newFile", fileInput.files[0]);

    try {
        const resp = await fetch('/upload', {
            method: 'POST',
            body: formData
        });
        if (resp.ok) {
            alert("Upload complete");
            fileUploaded = true;
            fileName = fileInput.files[0].name;
        } else {
            alert("Upload failed");
        }
    } catch (error) {
        console.error("Error uploading file:", error);
    }
}

async function postTask() {
    if (!fileUploaded || !fileName) {
        alert("Upload file first");
        return;
    }

    const taskData = {
        file_name: fileName,
        output_format: "png"
    };

    try {
        const resp = await fetch('/task', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(taskData)
        });
        if(resp.ok) {
            alert("Task Queued");
            const data = await resp.json();
            outputFileName = data.output_file;
            outputFileNameP.textContent = `Output File: ${data.output_file}`;
        } else {
            alert("Failed to post task");
        }
    } catch (error) {
        console.error("Error posting task: ", error);
    }

    fileUploaded = false;
    fileName = '';
    fileInput.value = '';
}

async function downloadFile() {
    if(outputFileName == "") {
        alert("No file to download");
        return;
    }
    try {
        const resp = await fetch(`download/${encodeURIComponent(outputFileName)}`, {
            method: 'GET',
        });
        if (!resp.ok) {
            const errorText = await resp.text();
            console.error(`Download failed (${resp.status}):`, errorText);
            alert(`Failed to download file: ${errorText}`);
            return;
        }

        // Convert the successful response into a binary Blob
        const blob = await resp.blob();

        // Create a temporary object URL for the blob
        const downloadUrl = window.URL.createObjectURL(blob);

        // Create a hidden <a> tag to trigger the browser's save dialog
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = outputFileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        // Clean up memory
        window.URL.revokeObjectURL(downloadUrl);
    } catch (error) {
        console.error('Network error during file download:', error);
        alert('An error occurred while downloading the file.');
    }
}

async function clearFields() {
    fileUploaded = false;
    fileName = "";
    outputFileName = "";
    fileInput.value = '';
    outputFileNameP.textContent = "";
}