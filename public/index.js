console.log('init');

function removeFromDb(title){
  fetch(`/delete?title=${title}`, {method: "DELETE"}).then(res =>{
      if (res.status == 200){
          alert("Entry deleted");
          window.location.pathname = "/"
      }
  })
}

function updateDb(title) {
  const newTitle = document.getElementById("Add_Item").value;

  fetch(`/update?oldTitle=${title}&newTitle=${newTitle}`, {method: "PUT"}).then(res =>{
      if (res.status == 200){
          alert("Database updated")
          window.location.pathname = "/"
      }
  })
}