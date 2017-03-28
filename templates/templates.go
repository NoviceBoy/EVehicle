package templates

const refreshHeader = `
<!DOCTYPE html>
<html lang="en">
<head>
  <title>Electric Vehicle Dashboard</title>
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="20">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
</head>
<body>

<nav class="navbar navbar-inverse">
  <div class="container-fluid">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">Dashboard</a>
    </div>
    <ul class="nav navbar-nav">
      <li><a href="/">Home</a></li>
      <li><a href="/addUser">Register Vehicle</a></li>
      <li><a href="/showAllUsers">All Users</a></li>
      <li><a href="/notify">Notify</a></li>
    </ul>
  </div>
</nav>

<div class="container">
`
const header = `
<!DOCTYPE html>
<html lang="en">
<head>
  <title>Electric Vehicle Dashboard</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
</head>
<body>

<nav class="navbar navbar-inverse">
  <div class="container-fluid">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">Dashboard</a>
    </div>
    <ul class="nav navbar-nav">
      <li><a href="/">Home</a></li>
      <li><a href="/addUser">Register Vehicle</a></li>
      <li><a href="/showAllUsers">All Users</a></li>
      <li><a href="/notify">Notify</a></li>
    </ul>
  </div>
</nav>

<div class="container">
`

const footer = `
    </div>
    </body>
  </html>`

const ResultsPage = header + ` 
      <h3>Owner of {{.LicensePlate}} has been notified </h3>
      <table class="table table-hover table-bordered">
          <thead>
              <tr>
                  <th>License Plate</th>
                  <th>Name</th>
                  <th>Email</th>
              </tr>
          </thead>
          <tbody>
            <tr>
              <td>{{.LicensePlate}}</td>
              <td>{{.Name}}</td>
              <td>{{.Email}}</td>
            </tr>
          </tbody>
      </table> 
` + footer

const ErrorPage = header + ` 
      <p><b>{{html .}}</b></p>
      <p><a href="/">Return to main page</a></p>
` + footer

const NotifyPage = header + `
        <h4>License Plate Notifier</h4><br>
        <form class="form-horizontal" action="/result" method="post">
          <div class="form-group">
            <div class="col-sm-6">
              <input type="text" name="licensePlate" class="form-control" placeholder="Enter License Plate">
            </div>
          </div>
          <div class="form-group">
            <div class="col-sm-6">
              <button type="submit" class="btn btn-success">Notify</button>
            </div>
          </div>
          </div>   
        <br>
` + footer

const AddUserPage = header + `
        <h3>Vehicle Registration</h3><br>
        <form class="form-horizontal" action="/addedUser" method="post" >
          <div class="form-group">
            <label class="col-sm-2" for="licensePlate">License Plate:</label>
            <div class="col-sm-8">
              <input type="text" class="form-control" id="licensePlate" placeholder="License plate" name="licensePlate" required/>
            </div>
          </div>
          <div class="form-group">
            <label class="col-sm-2" for="name">Name:</label>
            <div class="col-sm-8">
              <input type="text" class="form-control" id="name" placeholder="Name" name="name" required/>
            </div>
          </div>
          <div class="form-group">
            <label class="col-sm-2" for="email">Email:</label>
            <div class="col-sm-8">
              <input type="text" class="form-control" id="email" placeholder="Email" name="email" required/>
            </div>
          </div>
          <button type="submit" class="btn btn-success">Submit</button>
        </form>
` + footer

const AllResultsPage = header + ` 
      <table class="table table-hover table-bordered">
          <thead>
              <tr>
                  <th>License Plate</th>
                  <th>Name</th>
                  <th>Email</th>
              </tr>
          </thead>
          <tbody>
            {{ range . }}
            <tr>
              <td>{{.LicensePlate}}</td>
              <td>{{.Name}}</td>
              <td>{{.Email}}</td>
            </tr>
            {{ end }}
          </tbody>
      </table> 
` + footer

const DashBoardPage = refreshHeader + ` 
    <div class="col-sm-8">
    <h4>Station dashboard</h4> 
      <table class="table table-hover table-bordered">
          <thead>
              <tr>
                  <th>Charging Station</th>
                  <th>Current User</th>
              </tr>
          </thead>
          <tbody>
            {{ range .Stations }}
            <tr>
              <td>{{.Id}}</td>
              {{ if .User.Name }}
              <td>{{.User.Name }}</td>
              {{ else }}
              <td><div class="btn btn-success"> FREE </div></td>
              {{ end }}
            </tr>
            {{ end }}
          </tbody>
      </table> 
    </div>
    <div class="col-sm-4">
    <h4> User queue </h4>
      <table class="table table-hover table-bordered">
          <thead>
              <tr>
                  <th>Queue</th>
                  <th>Status</th>
              </tr>
          </thead>
          <tbody>
            {{ range .Queue }}
            <tr>
              <td>{{.Name }}</td>
              <td>{{.Status }}</td>
            </tr>
            {{ end }}
          </tbody>
      </table> 
      </div>
  ` + footer

const modelButton = `<button type="button" class="btn btn-success" data-toggle="modal" data-target="#user">Add to Queue</button>
        <div class="modal fade" id="user" role="dialog">
          <div class="modal-dialog">
            <div class="modal-content">
              <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal">&times;</button>
                <h4 class="modal-title">Add User to Queue</h4>
              </div>
              <div class="modal-body">
                <form class="form-horizontal" action="/" method="post" >
                  <div class="form-group">
                    <label class="col-sm-3" for="licensePlate">License Plate:</label>
                    <div class="col-sm-9">
                      <input type="text" class="form-control" id="licensePlate" placeholder="License plate" name="licensePlate" required/>
                    </div>
                  </div>
                  <div class="form-group">
                    <label class="col-sm-3" for="name">Name:</label>
                    <div class="col-sm-9">
                      <input type="text" class="form-control" id="name" placeholder="Name" name="name" required/>
                    </div>
                  </div>
                  <button type="submit" class="btn btn-success">Submit</button>
                </form>
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
              </div>
            </div>  
          </div>
        </div>`
