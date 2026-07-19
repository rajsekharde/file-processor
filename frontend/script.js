var fileUploaded = false;
var fileName = "";

const fileInput = document.getElementById("fileInput");
const uploadBtn = document.getElementById("fileUploadButton");
const postTaskBtn = document.getElementById("postTaskButton");

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
        convert_to: "jpg"
    };

    try {
        const resp =  await fetch('/task', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(taskData)
        });
        if(resp.ok) {
            alert("Task Queued");
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