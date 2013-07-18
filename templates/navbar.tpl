<div class="navbar navbar-fixed-top">
	<div class="navbar-inner">
		<div class="container">
			<a href="/" class="brand">
				Go Daily
				<sub class="version">{{.T.AppVer}}</sub>
			</a>
			<ul class="nav">
				<li {{if .T.IsHome}}class="active"{{end}}><a href="/">广场</a></li>
				<li {{if .T.IsExer}}class="active"{{end}}><a href="/exercise">每日一练</a></li>
				<li {{if .T.IsAbout}}class="active"{{end}}><a href="/about">关于</a></li>
			</ul>
		    <div class="pull-right">
		        <ul class="nav">
		            <li>
		        	{{if IsLogedIn}}
		        		<a href="{{.App.BasePath}}{{GetLoginUserName}}">
		        			<img src="{{GetLoginUserAvatar}}?s=16"/> {{GetLoginUserName}}</a>
		        		<a href="{{.App.BasePath}}logout">退出</a>
		        	{{else}}
		        		<a href="/login">登录</a>
		        	{{end}}
		            </li>	

		            <li>
		            	<a href="/register">注册</a>
		            </li>
		        </ul>
		        </div>
		    </div>
		</div>
	</div>
</div>