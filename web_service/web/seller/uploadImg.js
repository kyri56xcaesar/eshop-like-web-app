let queuedImagesArray = [],
inputDiv = document.querySelector(".input-div"),
input = document.querySelector(".input-div input"),
imageDepiction = document.getElementById("imageDepiction");
deletedImages = [];

var productForm = document.getElementById("product-form");

//console.log("Element");
// console.log(input);

input.addEventListener("change", () => {
  const files = input.files;
  console.log(files);

  for (let i = 0; i < files.length; i++) {
    queuedImagesArray.push(files[i]);
  }
  displayQueuedImages();
});

inputDiv.addEventListener("drop", (e) => {
  e.preventDefault();
  const files = e.dataTransfer.files
  // console.log(files);

  for (let i = 0; i < files.length; i++){
    if (!files[i].type.match("image")) continue;

    if (queuedImagesArray.every(image => image.name !== files[i].name))
      queuedImagesArray.push(files[i]);
  }
  displayQueuedImages();
});

function displayQueuedImages() {
  let images = "";
  queuedImagesArray.forEach((image, index) => {
    images += `<div class="image">
                <img id="product-image" src="${URL.createObjectURL(image)}" alt="image">
                <span onclick="deleteQueuedImage(${index})">&times;</span>
               </div>`
  });
  imageDepiction.innerHTML = images;

  
}

function deleteQueuedImage(index) {
  queuedImagesArray.splice(index, 1);
  displayQueuedImages();
}

productForm.addEventListener("submit", (e) => {
  e.preventDefault();
  sendImagesToServer();
});


function sendImagesToServer() {
  const formData = new FormData(productForm);

  queuedImagesArray.forEach((image, index) => {
    formData.append(`file[${index}]`, image);
  });

  
}
