async function uploadFile() {
    const fileInput = document.getElementById("fileInput");
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
        if (response.ok) {
            alert("Upload complete");
        } else {
            alert("Upload failed");
        }
    } catch (error) {
        console.error("Error uploading file:", error);
    }
}