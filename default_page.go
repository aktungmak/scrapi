package scrapi

var BUILD_TIME = "unknown"

const DEFAULT_PAGE = `
<html>
<head>
<title>scrapi</title>
</head>
<body>
<h2>scrapi</h2>
<p>
Browse the API starting from the 
<a href="{{.R.Root}}">service root</a>.</p>
<p>
<b>Capture time:</b> {{.R.Date}}</p>
<p><b>Host:</b> {{.R.Host}}</p>
<p><b>Notes:</b> {{.R.Notes}}</p>

<hr/>
<footer>
scrapi built on {{.D}}. 
questions to <a href="mailto:joshua.green@ericsson.com">joshua.green@ericsson.com</a>
</footer>
</body>
</html>
`
