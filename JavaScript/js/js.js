
function time_startTime() {
    var today = new Date();
    var h = today.getHours();
    var m = today.getMinutes();
    var s = today.getSeconds();
    m = time_checkTime(m);
    s = time_checkTime(s);
    document.getElementById('time').innerHTML =
    h + ":" + m + ":" + s;
    var t = setTimeout(time_startTime, 500);
}
function time_checkTime(i) {
    if (i < 10) {i = "0" + i};
    return i;
}

function Submit_form(){
    var hoten = document.getElementById("hoten").value;
    var tmp_gioitinh = document.getElementById("gioitinh");
    var gioitinh = tmp_gioitinh.options[tmp_gioitinh.selectedIndex].value;
    var ngaysinh = document.getElementById("ngaysinh").value;
    var object_ngaysinh = new Date(ngaysinh);
    var object_ngay = object_ngaysinh.getDate();
    var object_thang = object_ngaysinh.getMonth()+1;
    var object_nam = object_ngaysinh.getFullYear();
    document.getElementById("replace_hoten").innerHTML = hoten;
    document.getElementById("replace_ngaysinh").innerHTML = object_ngay+"/"+object_thang+"/"+object_nam;
    document.getElementById("replace_gioitinh").innerHTML = gioitinh;
}
function hide_img(){
    var x = document.getElementById("img-change");
    if (x.style.display === "none") {
      x.style.display = "block";
      document.getElementById("show").style.display = "none";
    } else {
      x.style.display = "none";
      document.getElementById("show").style.display = "block";
    }
}
