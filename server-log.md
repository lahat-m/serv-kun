mis-app  | 2026-08-07T10:57:38.772Z  WARN 1 --- [mis] [nio-8080-exec-9] com.zaxxer.hikari.pool.PoolBase          : HikariPool-1 - Failed to validate connection org.postgresql.jdbc.PgConnection@22a0c6c (This connection has been closed.). Possibly consider using a shorter maxLifetime value.
mis-app  | 2026-08-07T10:58:08.770Z  WARN 1 --- [mis] [nio-8080-exec-9] o.s.b.j.h.DataSourceHealthIndicator      : DataSource health check failed
mis-app  |
mis-app  | org.springframework.jdbc.CannotGetJdbcConnectionException: Failed to obtain JDBC Connection
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:84) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.core.JdbcTemplate.execute(JdbcTemplate.java:362) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.getProduct(DataSourceHealthIndicator.java:125) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doDataSourceHealthCheck(DataSourceHealthIndicator.java:108) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doHealthCheck(DataSourceHealthIndicator.java:102) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.AbstractHealthIndicator.health(AbstractHealthIndicator.java:80) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.HealthIndicator.health(HealthIndicator.java:37) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:148) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:126) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptorAndLogIfSlow(HealthEndpointSupport.java:172) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:145) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getAggregateDescriptor(HealthEndpointSupport.java:157) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:141) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:110) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:82) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:85) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:73) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.util.ReflectionUtils.invokeMethod(ReflectionUtils.java:281) ~[spring-core-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.actuate.endpoint.invoke.reflect.ReflectiveOperationInvoker.invoke(ReflectiveOperationInvoker.java:76) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.actuate.endpoint.annotation.AbstractDiscoveredOperation.invoke(AbstractDiscoveredOperation.java:62) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$ServletWebOperationAdapter.handle(AbstractWebMvcEndpointHandlerMapping.java:347) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$OperationHandler.handle(AbstractWebMvcEndpointHandlerMapping.java:476) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.doInvoke(InvocableHandlerMethod.java:252) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.invokeForRequest(InvocableHandlerMethod.java:184) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.ServletInvocableHandlerMethod.invokeAndHandle(ServletInvocableHandlerMethod.java:117) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.invokeHandlerMethod(RequestMappingHandlerAdapter.java:934) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.handleInternal(RequestMappingHandlerAdapter.java:853) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.AbstractHandlerMethodAdapter.handle(AbstractHandlerMethodAdapter.java:86) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:963) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:866) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1000) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.doGet(FrameworkServlet.java:892) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:622) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.service(FrameworkServlet.java:874) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:710) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:128) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53) ~[tomcat-embed-websocket-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.servlet.resource.ResourceUrlEncodingFilter.doFilter(ResourceUrlEncodingFilter.java:66) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:108) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.FilterChainProxy.lambda$doFilterInternal$2(FilterChainProxy.java:235) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$FilterObservation$SimpleFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:491) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:353) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator.lambda$wrapSecured$0(ObservationFilterChainDecorator.java:84) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:131) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.intercept.AuthorizationFilter.doFilter(AuthorizationFilter.java:101) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:126) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:120) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:132) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:86) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.AnonymousAuthenticationFilter.doFilter(AnonymousAuthenticationFilter.java:100) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.servletapi.SecurityContextHolderAwareRequestFilter.doFilter(SecurityContextHolderAwareRequestFilter.java:181) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.savedrequest.RequestCacheAwareFilter.doFilter(RequestCacheAwareFilter.java:63) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.authentication.BearerTokenAuthenticationFilter.doFilterInternal(BearerTokenAuthenticationFilter.java:174) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.OAuth2ProtectedResourceMetadataFilter.doFilterInternal(OAuth2ProtectedResourceMetadataFilter.java:98) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:110) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:96) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CorsFilter.doFilterInternal(CorsFilter.java:91) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doHeadersAfter(HeaderWriterFilter.java:90) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doFilterInternal(HeaderWriterFilter.java:75) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:82) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:69) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.request.async.WebAsyncManagerIntegrationFilter.doFilterInternal(WebAsyncManagerIntegrationFilter.java:62) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.DisableEncodeUrlFilter.doFilterInternal(DisableEncodeUrlFilter.java:42) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$0(ObservationFilterChainDecorator.java:336) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:227) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilterInternal(FilterChainProxy.java:237) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilter(FilterChainProxy.java:195) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.ServletRequestPathFilter.doFilter(ServletRequestPathFilter.java:52) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter.doFilter(CompositeFilter.java:74) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.config.annotation.web.configuration.WebSecurityConfiguration$CompositeFilterChainProxy.doFilter(WebSecurityConfiguration.java:317) ~[spring-security-config-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.invokeDelegate(DelegatingFilterProxy.java:355) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.doFilter(DelegatingFilterProxy.java:272) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.ServerHttpObservationFilter.doFilterInternal(ServerHttpObservationFilter.java:110) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:199) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:165) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:77) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:492) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:113) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:83) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:72) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:341) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:397) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:63) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:1272) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1801) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:52) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:946) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:480) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:57) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.sql.SQLTransientConnectionException: HikariPool-1 - Connection is not available, request timed out after 30000ms (total=0, active=0, idle=0, waiting=0)
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createTimeoutException(HikariPool.java:714) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:184) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:142) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.HikariDataSource.getConnection(HikariDataSource.java:127) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.fetchConnection(DataSourceUtils.java:160) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.doGetConnection(DataSourceUtils.java:118) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:81) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      ... 163 common frames omitted
mis-app  | Caused by: org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      ... 1 common frames omitted
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T10:58:08.788Z  WARN 1 --- [mis] [nio-8080-exec-9] o.s.b.h.a.e.HealthEndpointSupport        : Health contributor org.springframework.boot.jdbc.health.DataSourceHealthIndicator (db) took 30013ms to respond
mis-app  | 2026-08-07T10:58:27.718Z  WARN 1 --- [mis] [onnection-adder] com.zaxxer.hikari.pool.PoolBase          : HikariPool-1 - Pool is empty, failed to create/setup connection (1d23b94c-accc-4f69-b0a6-f25c54578ca6)
mis-app  |
mis-app  | org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T10:58:48.817Z  WARN 1 --- [mis] [nio-8080-exec-3] o.s.b.j.h.DataSourceHealthIndicator      : DataSource health check failed
mis-app  |
mis-app  | org.springframework.jdbc.CannotGetJdbcConnectionException: Failed to obtain JDBC Connection
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:84) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.core.JdbcTemplate.execute(JdbcTemplate.java:362) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.getProduct(DataSourceHealthIndicator.java:125) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doDataSourceHealthCheck(DataSourceHealthIndicator.java:108) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doHealthCheck(DataSourceHealthIndicator.java:102) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.AbstractHealthIndicator.health(AbstractHealthIndicator.java:80) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.HealthIndicator.health(HealthIndicator.java:37) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:148) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:126) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptorAndLogIfSlow(HealthEndpointSupport.java:172) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:145) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getAggregateDescriptor(HealthEndpointSupport.java:157) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:141) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:110) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:82) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:85) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:73) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.util.ReflectionUtils.invokeMethod(ReflectionUtils.java:281) ~[spring-core-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.actuate.endpoint.invoke.reflect.ReflectiveOperationInvoker.invoke(ReflectiveOperationInvoker.java:76) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.actuate.endpoint.annotation.AbstractDiscoveredOperation.invoke(AbstractDiscoveredOperation.java:62) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$ServletWebOperationAdapter.handle(AbstractWebMvcEndpointHandlerMapping.java:347) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$OperationHandler.handle(AbstractWebMvcEndpointHandlerMapping.java:476) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.doInvoke(InvocableHandlerMethod.java:252) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.invokeForRequest(InvocableHandlerMethod.java:184) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.ServletInvocableHandlerMethod.invokeAndHandle(ServletInvocableHandlerMethod.java:117) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.invokeHandlerMethod(RequestMappingHandlerAdapter.java:934) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.handleInternal(RequestMappingHandlerAdapter.java:853) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.AbstractHandlerMethodAdapter.handle(AbstractHandlerMethodAdapter.java:86) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:963) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:866) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1000) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.doGet(FrameworkServlet.java:892) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:622) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.service(FrameworkServlet.java:874) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:710) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:128) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53) ~[tomcat-embed-websocket-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.servlet.resource.ResourceUrlEncodingFilter.doFilter(ResourceUrlEncodingFilter.java:66) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:108) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.FilterChainProxy.lambda$doFilterInternal$2(FilterChainProxy.java:235) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$FilterObservation$SimpleFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:491) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:353) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator.lambda$wrapSecured$0(ObservationFilterChainDecorator.java:84) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:131) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.intercept.AuthorizationFilter.doFilter(AuthorizationFilter.java:101) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:126) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:120) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:132) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:86) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.AnonymousAuthenticationFilter.doFilter(AnonymousAuthenticationFilter.java:100) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.servletapi.SecurityContextHolderAwareRequestFilter.doFilter(SecurityContextHolderAwareRequestFilter.java:181) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.savedrequest.RequestCacheAwareFilter.doFilter(RequestCacheAwareFilter.java:63) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.authentication.BearerTokenAuthenticationFilter.doFilterInternal(BearerTokenAuthenticationFilter.java:174) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.OAuth2ProtectedResourceMetadataFilter.doFilterInternal(OAuth2ProtectedResourceMetadataFilter.java:98) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:110) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:96) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CorsFilter.doFilterInternal(CorsFilter.java:91) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doHeadersAfter(HeaderWriterFilter.java:90) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doFilterInternal(HeaderWriterFilter.java:75) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:82) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:69) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.request.async.WebAsyncManagerIntegrationFilter.doFilterInternal(WebAsyncManagerIntegrationFilter.java:62) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.DisableEncodeUrlFilter.doFilterInternal(DisableEncodeUrlFilter.java:42) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$0(ObservationFilterChainDecorator.java:336) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:227) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilterInternal(FilterChainProxy.java:237) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilter(FilterChainProxy.java:195) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.ServletRequestPathFilter.doFilter(ServletRequestPathFilter.java:52) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter.doFilter(CompositeFilter.java:74) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.config.annotation.web.configuration.WebSecurityConfiguration$CompositeFilterChainProxy.doFilter(WebSecurityConfiguration.java:317) ~[spring-security-config-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.invokeDelegate(DelegatingFilterProxy.java:355) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.doFilter(DelegatingFilterProxy.java:272) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.ServerHttpObservationFilter.doFilterInternal(ServerHttpObservationFilter.java:110) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:199) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:165) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:77) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:492) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:113) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:83) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:72) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:341) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:397) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:63) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:1272) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1801) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:52) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:946) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:480) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:57) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.sql.SQLTransientConnectionException: HikariPool-1 - Connection is not available, request timed out after 30000ms (total=0, active=0, idle=0, waiting=0)
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createTimeoutException(HikariPool.java:714) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:184) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:142) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.HikariDataSource.getConnection(HikariDataSource.java:127) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.fetchConnection(DataSourceUtils.java:160) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.doGetConnection(DataSourceUtils.java:118) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:81) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      ... 163 common frames omitted
mis-app  | Caused by: org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      ... 1 common frames omitted
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T10:58:48.819Z  WARN 1 --- [mis] [nio-8080-exec-3] o.s.b.h.a.e.HealthEndpointSupport        : Health contributor org.springframework.boot.jdbc.health.DataSourceHealthIndicator (db) took 30002ms to respond
mis-app  | 2026-08-07T10:59:28.862Z  WARN 1 --- [mis] [nio-8080-exec-8] o.s.b.j.h.DataSourceHealthIndicator      : DataSource health check failed
mis-app  |
mis-app  | org.springframework.jdbc.CannotGetJdbcConnectionException: Failed to obtain JDBC Connection
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:84) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.core.JdbcTemplate.execute(JdbcTemplate.java:362) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.getProduct(DataSourceHealthIndicator.java:125) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doDataSourceHealthCheck(DataSourceHealthIndicator.java:108) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doHealthCheck(DataSourceHealthIndicator.java:102) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.AbstractHealthIndicator.health(AbstractHealthIndicator.java:80) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.HealthIndicator.health(HealthIndicator.java:37) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:148) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:126) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptorAndLogIfSlow(HealthEndpointSupport.java:172) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:145) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getAggregateDescriptor(HealthEndpointSupport.java:157) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:141) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:110) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:82) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:85) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:73) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.util.ReflectionUtils.invokeMethod(ReflectionUtils.java:281) ~[spring-core-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.actuate.endpoint.invoke.reflect.ReflectiveOperationInvoker.invoke(ReflectiveOperationInvoker.java:76) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.actuate.endpoint.annotation.AbstractDiscoveredOperation.invoke(AbstractDiscoveredOperation.java:62) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$ServletWebOperationAdapter.handle(AbstractWebMvcEndpointHandlerMapping.java:347) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$OperationHandler.handle(AbstractWebMvcEndpointHandlerMapping.java:476) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.doInvoke(InvocableHandlerMethod.java:252) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.invokeForRequest(InvocableHandlerMethod.java:184) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.ServletInvocableHandlerMethod.invokeAndHandle(ServletInvocableHandlerMethod.java:117) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.invokeHandlerMethod(RequestMappingHandlerAdapter.java:934) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.handleInternal(RequestMappingHandlerAdapter.java:853) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.AbstractHandlerMethodAdapter.handle(AbstractHandlerMethodAdapter.java:86) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:963) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:866) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1000) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.doGet(FrameworkServlet.java:892) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:622) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.service(FrameworkServlet.java:874) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:710) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:128) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53) ~[tomcat-embed-websocket-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.servlet.resource.ResourceUrlEncodingFilter.doFilter(ResourceUrlEncodingFilter.java:66) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:108) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.FilterChainProxy.lambda$doFilterInternal$2(FilterChainProxy.java:235) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$FilterObservation$SimpleFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:491) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:353) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator.lambda$wrapSecured$0(ObservationFilterChainDecorator.java:84) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:131) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.intercept.AuthorizationFilter.doFilter(AuthorizationFilter.java:101) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:126) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:120) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:132) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:86) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.AnonymousAuthenticationFilter.doFilter(AnonymousAuthenticationFilter.java:100) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.servletapi.SecurityContextHolderAwareRequestFilter.doFilter(SecurityContextHolderAwareRequestFilter.java:181) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.savedrequest.RequestCacheAwareFilter.doFilter(RequestCacheAwareFilter.java:63) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.authentication.BearerTokenAuthenticationFilter.doFilterInternal(BearerTokenAuthenticationFilter.java:174) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.OAuth2ProtectedResourceMetadataFilter.doFilterInternal(OAuth2ProtectedResourceMetadataFilter.java:98) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:110) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:96) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CorsFilter.doFilterInternal(CorsFilter.java:91) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doHeadersAfter(HeaderWriterFilter.java:90) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doFilterInternal(HeaderWriterFilter.java:75) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:82) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:69) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.request.async.WebAsyncManagerIntegrationFilter.doFilterInternal(WebAsyncManagerIntegrationFilter.java:62) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.DisableEncodeUrlFilter.doFilterInternal(DisableEncodeUrlFilter.java:42) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$0(ObservationFilterChainDecorator.java:336) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:227) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilterInternal(FilterChainProxy.java:237) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilter(FilterChainProxy.java:195) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.ServletRequestPathFilter.doFilter(ServletRequestPathFilter.java:52) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter.doFilter(CompositeFilter.java:74) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.config.annotation.web.configuration.WebSecurityConfiguration$CompositeFilterChainProxy.doFilter(WebSecurityConfiguration.java:317) ~[spring-security-config-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.invokeDelegate(DelegatingFilterProxy.java:355) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.doFilter(DelegatingFilterProxy.java:272) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.ServerHttpObservationFilter.doFilterInternal(ServerHttpObservationFilter.java:110) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:199) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:165) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:77) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:492) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:113) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:83) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:72) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:341) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:397) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:63) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:1272) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1801) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:52) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:946) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:480) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:57) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.sql.SQLTransientConnectionException: HikariPool-1 - Connection is not available, request timed out after 30000ms (total=0, active=0, idle=0, waiting=0)
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createTimeoutException(HikariPool.java:714) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:184) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:142) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.HikariDataSource.getConnection(HikariDataSource.java:127) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.fetchConnection(DataSourceUtils.java:160) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.doGetConnection(DataSourceUtils.java:118) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:81) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      ... 163 common frames omitted
mis-app  | Caused by: org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      ... 1 common frames omitted
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T10:59:28.865Z  WARN 1 --- [mis] [nio-8080-exec-8] o.s.b.h.a.e.HealthEndpointSupport        : Health contributor org.springframework.boot.jdbc.health.DataSourceHealthIndicator (db) took 30003ms to respond
mis-app  | 2026-08-07T10:59:32.740Z  WARN 1 --- [mis] [onnection-adder] com.zaxxer.hikari.pool.PoolBase          : HikariPool-1 - Pool is empty, failed to create/setup connection (00b2b0a4-a3a6-4c56-bae3-8f45d5059bec)
mis-app  |
mis-app  | org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T11:00:08.930Z  WARN 1 --- [mis] [nio-8080-exec-6] o.s.b.j.h.DataSourceHealthIndicator      : DataSource health check failed
mis-app  |
mis-app  | org.springframework.jdbc.CannotGetJdbcConnectionException: Failed to obtain JDBC Connection
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:84) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.core.JdbcTemplate.execute(JdbcTemplate.java:362) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.getProduct(DataSourceHealthIndicator.java:125) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doDataSourceHealthCheck(DataSourceHealthIndicator.java:108) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.jdbc.health.DataSourceHealthIndicator.doHealthCheck(DataSourceHealthIndicator.java:102) ~[spring-boot-jdbc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.AbstractHealthIndicator.health(AbstractHealthIndicator.java:80) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.contributor.HealthIndicator.health(HealthIndicator.java:37) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:148) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.Contributor$Blocking.getDescriptor(Contributor.java:126) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptorAndLogIfSlow(HealthEndpointSupport.java:172) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:145) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getAggregateDescriptor(HealthEndpointSupport.java:157) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getDescriptor(HealthEndpointSupport.java:141) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:110) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointSupport.getResult(HealthEndpointSupport.java:82) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:85) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.health.actuate.endpoint.HealthEndpointWebExtension.health(HealthEndpointWebExtension.java:73) ~[spring-boot-health-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.util.ReflectionUtils.invokeMethod(ReflectionUtils.java:281) ~[spring-core-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.boot.actuate.endpoint.invoke.reflect.ReflectiveOperationInvoker.invoke(ReflectiveOperationInvoker.java:76) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.actuate.endpoint.annotation.AbstractDiscoveredOperation.invoke(AbstractDiscoveredOperation.java:62) ~[spring-boot-actuator-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$ServletWebOperationAdapter.handle(AbstractWebMvcEndpointHandlerMapping.java:347) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at org.springframework.boot.webmvc.actuate.endpoint.web.AbstractWebMvcEndpointHandlerMapping$OperationHandler.handle(AbstractWebMvcEndpointHandlerMapping.java:476) ~[spring-boot-webmvc-4.1.0.jar:4.1.0]
mis-app  |      at java.base/jdk.internal.reflect.DirectMethodHandleAccessor.invoke(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.lang.reflect.Method.invoke(Unknown Source) ~[na:na]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.doInvoke(InvocableHandlerMethod.java:252) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.method.support.InvocableHandlerMethod.invokeForRequest(InvocableHandlerMethod.java:184) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.ServletInvocableHandlerMethod.invokeAndHandle(ServletInvocableHandlerMethod.java:117) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.invokeHandlerMethod(RequestMappingHandlerAdapter.java:934) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerAdapter.handleInternal(RequestMappingHandlerAdapter.java:853) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.mvc.method.AbstractHandlerMethodAdapter.handle(AbstractHandlerMethodAdapter.java:86) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:963) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:866) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1000) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.doGet(FrameworkServlet.java:892) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:622) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.springframework.web.servlet.FrameworkServlet.service(FrameworkServlet.java:874) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at jakarta.servlet.http.HttpServlet.service(HttpServlet.java:710) ~[tomcat-embed-core-11.0.22.jar:6.1]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:128) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.websocket.server.WsFilter.doFilter(WsFilter.java:53) ~[tomcat-embed-websocket-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.servlet.resource.ResourceUrlEncodingFilter.doFilter(ResourceUrlEncodingFilter.java:66) ~[spring-webmvc-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:108) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.FilterChainProxy.lambda$doFilterInternal$2(FilterChainProxy.java:235) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$FilterObservation$SimpleFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:491) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$1(ObservationFilterChainDecorator.java:353) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator.lambda$wrapSecured$0(ObservationFilterChainDecorator.java:84) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:131) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.intercept.AuthorizationFilter.doFilter(AuthorizationFilter.java:101) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:126) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.access.ExceptionTranslationFilter.doFilter(ExceptionTranslationFilter.java:120) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:132) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.SessionManagementFilter.doFilter(SessionManagementFilter.java:86) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.AnonymousAuthenticationFilter.doFilter(AnonymousAuthenticationFilter.java:100) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.servletapi.SecurityContextHolderAwareRequestFilter.doFilter(SecurityContextHolderAwareRequestFilter.java:181) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.savedrequest.RequestCacheAwareFilter.doFilter(RequestCacheAwareFilter.java:63) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.authentication.BearerTokenAuthenticationFilter.doFilterInternal(BearerTokenAuthenticationFilter.java:174) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.oauth2.server.resource.web.OAuth2ProtectedResourceMetadataFilter.doFilterInternal(OAuth2ProtectedResourceMetadataFilter.java:98) ~[spring-security-oauth2-resource-server-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:110) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.authentication.logout.LogoutFilter.doFilter(LogoutFilter.java:96) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CorsFilter.doFilterInternal(CorsFilter.java:91) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doHeadersAfter(HeaderWriterFilter.java:90) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.header.HeaderWriterFilter.doFilterInternal(HeaderWriterFilter.java:75) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:82) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.SecurityContextHolderFilter.doFilter(SecurityContextHolderFilter.java:69) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.context.request.async.WebAsyncManagerIntegrationFilter.doFilterInternal(WebAsyncManagerIntegrationFilter.java:62) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:230) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.session.DisableEncodeUrlFilter.doFilterInternal(DisableEncodeUrlFilter.java:42) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.wrapFilter(ObservationFilterChainDecorator.java:243) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$AroundFilterObservation$SimpleAroundFilterObservation.lambda$wrap$0(ObservationFilterChainDecorator.java:336) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$ObservationFilter.doFilter(ObservationFilterChainDecorator.java:227) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.ObservationFilterChainDecorator$VirtualFilterChain.doFilter(ObservationFilterChainDecorator.java:140) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilterInternal(FilterChainProxy.java:237) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.security.web.FilterChainProxy.doFilter(FilterChainProxy.java:195) ~[spring-security-web-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.ServletRequestPathFilter.doFilter(ServletRequestPathFilter.java:52) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter$VirtualFilterChain.doFilter(CompositeFilter.java:113) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.CompositeFilter.doFilter(CompositeFilter.java:74) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.security.config.annotation.web.configuration.WebSecurityConfiguration$CompositeFilterChainProxy.doFilter(WebSecurityConfiguration.java:317) ~[spring-security-config-7.1.0.jar:7.1.0]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.invokeDelegate(DelegatingFilterProxy.java:355) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.DelegatingFilterProxy.doFilter(DelegatingFilterProxy.java:272) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.RequestContextFilter.doFilterInternal(RequestContextFilter.java:100) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.FormContentFilter.doFilterInternal(FormContentFilter.java:93) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.ServerHttpObservationFilter.doFilterInternal(ServerHttpObservationFilter.java:110) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.CharacterEncodingFilter.doFilterInternal(CharacterEncodingFilter.java:199) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:116) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.springframework.web.filter.OncePerRequestFilter.doFilter(OncePerRequestFilter.java:101) ~[spring-web-7.0.8.jar:7.0.8]
mis-app  |      at org.apache.catalina.core.ApplicationFilterChain.doFilter(ApplicationFilterChain.java:107) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardWrapperValve.invoke(StandardWrapperValve.java:165) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardContextValve.invoke(StandardContextValve.java:77) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.authenticator.AuthenticatorBase.invoke(AuthenticatorBase.java:492) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardHostValve.invoke(StandardHostValve.java:113) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.valves.ErrorReportValve.invoke(ErrorReportValve.java:83) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.core.StandardEngineValve.invoke(StandardEngineValve.java:72) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.catalina.connector.CoyoteAdapter.service(CoyoteAdapter.java:341) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:397) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:63) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:1272) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1801) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:52) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:946) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:480) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:57) ~[tomcat-embed-core-11.0.22.jar:11.0.22]
mis-app  |      at java.base/java.lang.Thread.run(Unknown Source) ~[na:na]
mis-app  | Caused by: java.sql.SQLTransientConnectionException: HikariPool-1 - Connection is not available, request timed out after 30000ms (total=0, active=0, idle=0, waiting=0)
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createTimeoutException(HikariPool.java:714) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:184) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.getConnection(HikariPool.java:142) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.HikariDataSource.getConnection(HikariDataSource.java:127) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.fetchConnection(DataSourceUtils.java:160) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.doGetConnection(DataSourceUtils.java:118) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      at org.springframework.jdbc.datasource.DataSourceUtils.getConnection(DataSourceUtils.java:81) ~[spring-jdbc-7.0.8.jar:7.0.8]
mis-app  |      ... 163 common frames omitted
mis-app  | Caused by: org.postgresql.util.PSQLException: The connection attempt failed.
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:437) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.ConnectionFactory.openConnection(ConnectionFactory.java:52) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.jdbc.PgConnection.<init>(PgConnection.java:290) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.makeConnection(Driver.java:448) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.Driver.connect(Driver.java:298) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at com.zaxxer.hikari.util.DriverDataSource.getConnection(DriverDataSource.java:144) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newConnection(PoolBase.java:373) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.PoolBase.newPoolEntry(PoolBase.java:210) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool.createPoolEntry(HikariPool.java:488) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:752) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at com.zaxxer.hikari.pool.HikariPool$PoolEntryCreator.call(HikariPool.java:731) ~[HikariCP-7.0.2.jar:na]
mis-app  |      at java.base/java.util.concurrent.FutureTask.run(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(Unknown Source) ~[na:na]
mis-app  |      ... 1 common frames omitted
mis-app  | Caused by: java.net.UnknownHostException: postgres
mis-app  |      at java.base/sun.nio.ch.NioSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.SocksSocketImpl.connect(Unknown Source) ~[na:na]
mis-app  |      at java.base/java.net.Socket.connect(Unknown Source) ~[na:na]
mis-app  |      at org.postgresql.core.PGStream.createSocket(PGStream.java:271) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.PGStream.<init>(PGStream.java:132) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.tryConnect(ConnectionFactoryImpl.java:219) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      at org.postgresql.core.v3.ConnectionFactoryImpl.openConnectionImpl(ConnectionFactoryImpl.java:365) ~[postgresql-42.7.11.jar:42.7.11]
mis-app  |      ... 14 common frames omitted
mis-app  |
mis-app  | 2026-08-07T11:00:08.933Z  WARN 1 --- [mis] [nio-8080-exec-6] o.s.b.h.a.e.HealthEndpointSupport        : Health contributor org.springframework.boot.jdbc.health.DataSourceHealthIndicator (db) took 30003ms to respond
^C

lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ docker ps
permission denied while trying to connect to the docker API at unix:///var/run/docker.sock
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker ps
CONTAINER ID   IMAGE                        COMMAND                  CREATED         STATUS                    PORTS                                                                                      NAMES
c9c4022663d0   postgres:18-alpine           "docker-entrypoint.s…"   3 minutes ago   Up 3 minutes (healthy)    0.0.0.0:5434->5432/tcp, [::]:5434->5432/tcp                                                splm-mis-postgres
37fe0b6e5bf4   minio/minio:latest           "/usr/bin/docker-ent…"   3 minutes ago   Up 3 minutes (healthy)    0.0.0.0:9090->9000/tcp, [::]:9090->9000/tcp, 0.0.0.0:9091->9001/tcp, [::]:9091->9001/tcp   splm-mis-minio
cba9d34ebc54   mailhog/mailhog:latest       "MailHog"                3 minutes ago   Up 3 minutes              0.0.0.0:1025->1025/tcp, [::]:1025->1025/tcp, 0.0.0.0:8025->8025/tcp, [::]:8025->8025/tcp   splm-mis-mailhog
354d24e5c361   ghcr.io/lahat-m/mis:latest   "java -XX:+UseContai…"   11 hours ago    Up 11 hours (unhealthy)   8080/tcp                                                                                   mis-app
edf6776515fc   postgres:18-alpine           "/bin/sh /scripts/ba…"   13 hours ago    Up 13 hours               5432/tcp                                                                                   mis-db-backup
8fbd8cd80472   nginx:1.27-alpine            "/docker-entrypoint.…"   13 hours ago    Up 13 hours               0.0.0.0:80->80/tcp, [::]:80->80/tcp, 0.0.0.0:443->443/tcp, [::]:443->443/tcp               mis-nginx
554d96e44567   certbot/certbot:latest       "/bin/sh -c 'trap ex…"   13 hours ago    Up 13 hours               80/tcp, 443/tcp                                                                            mis-certbot
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker exec mis-app printenv | grep SPRING_DATASOURCE
SPRING_DATASOURCE_USERNAME=mis
SPRING_DATASOURCE_PASSWORD=JKLO5*DcctwNu
SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/mis
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker inspect mis-app --format='{{range $k,$v := .NetworkSettings.Networks}}{{$k}}{{println}}{{end}}'
splm-mis_internal

lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker exec -it mis-app getent hosts postgres
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker exec -it mis-app ping -c 1 postgres
OCI runtime exec failed: exec failed: unable to start container process: exec: "ping": executable file not found in $PATH
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker inspect splm-mis-postgres --format='{{range $k,$v := .NetworkSettings.Networks}}{{$k}}{{println}}{{end}}'
splm-mis_default

lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker inspect splm-mis-postgres --format='{{json .NetworkSettings.Networks.splm-mis_internal.Aliases}}'
template parsing error: template: :1: bad character U+002D '-'
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker compose ps
NAME                IMAGE                        COMMAND                  SERVICE    CREATED         STATUS                    PORTS
mis-app             ghcr.io/lahat-m/mis:latest   "java -XX:+UseContai…"   app        11 hours ago    Up 11 hours (unhealthy)   8080/tcp
mis-certbot         certbot/certbot:latest       "/bin/sh -c 'trap ex…"   certbot    13 hours ago    Up 13 hours               80/tcp, 443/tcp
mis-db-backup       postgres:18-alpine           "/bin/sh /scripts/ba…"   backup     13 hours ago    Up 13 hours               5432/tcp
mis-nginx           nginx:1.27-alpine            "/docker-entrypoint.…"   nginx      13 hours ago    Up 13 hours               0.0.0.0:80->80/tcp, [::]:80->80/tcp, 0.0.0.0:443->443/tcp, [::]:443->443/tcp
splm-mis-mailhog    mailhog/mailhog:latest       "MailHog"                mailhog    8 minutes ago   Up 8 minutes              0.0.0.0:1025->1025/tcp, [::]:1025->1025/tcp, 0.0.0.0:8025->8025/tcp, [::]:8025->8025/tcp
splm-mis-minio      minio/minio:latest           "/usr/bin/docker-ent…"   minio      8 minutes ago   Up 8 minutes (healthy)    0.0.0.0:9090->9000/tcp, [::]:9090->9000/tcp, 0.0.0.0:9091->9001/tcp, [::]:9091->9001/tcp
splm-mis-postgres   postgres:18-alpine           "docker-entrypoint.s…"   postgres   8 minutes ago   Up 8 minutes (healthy)    0.0.0.0:5434->5432/tcp, [::]:5434->5432/tcp
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker network ls
NETWORK ID     NAME                DRIVER    SCOPE
59b9f4b51fd1   bridge              bridge    local
34af6ff96a77   host                host      local
f54531475a97   none                null      local
2d42828ba75e   splm-mis_default    bridge    local
08cc1e7b3ba4   splm-mis_internal   bridge    local
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker compose down
sudo docker compose up -d
[+] down 5/5
✔ Container splm-mis-minio-init Removed                                                                                                                                              0.0s
✔ Container splm-mis-postgres   Removed                                                                                                                                              0.2s
✔ Container splm-mis-mailhog    Removed                                                                                                                                              0.2s
✔ Container splm-mis-minio      Removed                                                                                                                                              0.2s
! Network splm-mis_default      Resource is still in use                                                                                                                             0.0s
WARN[0000] Found orphan containers (mis-app, mis-db-backup, mis-nginx, mis-certbot) for this project. If you removed or renamed this service in your compose file, you can run this command with the --remove-orphans flag to clean it up.
[+] up 4/4
✔ Container splm-mis-postgres   Started                                                                                                                                              0.4s
✔ Container splm-mis-mailhog    Started                                                                                                                                              0.4s
✔ Container splm-mis-minio      Healthy                                                                                                                                              5.9s
✔ Container splm-mis-minio-init Started                                                                                                                                              6.0s
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker compose up -d --remove-orphans
[+] up 8/8
✔ Container splm-mis-postgres   Running                                                                                                                                              0.0s
✔ Container splm-mis-minio      Healthy                                                                                                                                             10.7s
✔ Container splm-mis-mailhog    Running                                                                                                                                              0.0s
✔ Container mis-nginx           Removed                                                                                                                                              0.2s
✔ Container mis-app             Removed                                                                                                                                             10.2s
✔ Container mis-certbot         Removed                                                                                                                                              0.2s
✔ Container mis-db-backup       Removed                                                                                                                                             10.2s
✔ Container splm-mis-minio-init Started                                                                                                                                              0.1s
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker compose up -d
[+] up 4/4
✔ Container splm-mis-postgres   Running                                                                                                                                              0.0s
✔ Container splm-mis-minio      Healthy                                                                                                                                              0.5s
✔ Container splm-mis-mailhog    Running                                                                                                                                              0.0s
✔ Container splm-mis-minio-init Started                                                                                                                                              0.1s
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docoker network ls
sudo: 'docoker': command not found
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker network ls
NETWORK ID     NAME                DRIVER    SCOPE
59b9f4b51fd1   bridge              bridge    local
34af6ff96a77   host                host      local
f54531475a97   none                null      local
2d42828ba75e   splm-mis_default    bridge    local
08cc1e7b3ba4   splm-mis_internal   bridge    local
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ sudo docker compose config
name: splm-mis
services:
mailhog:
container_name: splm-mis-mailhog
image: mailhog/mailhog:latest
networks:
default: null
ports:
- mode: ingress
target: 1025
published: "1025"
protocol: tcp
- mode: ingress
target: 8025
published: "8025"
protocol: tcp
restart: always
minio:
command:
- server
- /data
- --console-address
- :9001
container_name: splm-mis-minio
environment:
MINIO_ROOT_PASSWORD: minioadmin
MINIO_ROOT_USER: minioadmin
healthcheck:
test:
- CMD
- curl
- -f
- http://localhost:9000/minio/health/live
timeout: 10s
interval: 30s
retries: 3
start_period: 10s
image: minio/minio:latest
networks:
default: null
ports:
- mode: ingress
target: 9000
published: "9090"
protocol: tcp
- mode: ingress
target: 9001
published: "9091"
protocol: tcp
restart: unless-stopped
volumes:
- type: volume
source: splm-mis-minio-data
target: /data
volume: {}
minio-init:
container_name: splm-mis-minio-init
depends_on:
minio:
condition: service_healthy
required: true
entrypoint:
- /bin/sh
- -c
- |

          mc alias set local http://minio:9000 minioadmin minioadmin;
          mc mb --ignore-existing local/mis;
          echo '{"CORSRules":[{"AllowedOrigin":["*"],"AllowedMethod":["GET","PUT","HEAD"],"AllowedHeader":["*"],"MaxAgeSeconds":3600}]}' > /tmp/cors.json;
          mc cors set local/mis /tmp/cors.json || echo 'mc cors not supported by this MinIO version — presigned PUT uploads from the browser will fail CORS preflight until this is set some other way.';
          echo 'Bucket ready.';
    image: minio/mc:latest
    networks:
      default: null
    restart: 'no'
postgres:
container_name: splm-mis-postgres
environment:
POSTGRES_DB: mis
POSTGRES_PASSWORD: mis
POSTGRES_USER: mis
healthcheck:
test:
- CMD-SHELL
- pg_isready -U mis -d mis
timeout: 5s
interval: 10s
retries: 5
image: postgres:18-alpine
networks:
default: null
ports:
- mode: ingress
target: 5432
published: "5434"
protocol: tcp
restart: unless-stopped
volumes:
- type: volume
source: splm-mis-postgres-data
target: /var/lib/postgresql
volume: {}
networks:
default:
name: splm-mis_default
volumes:
splm-mis-minio-data:
name: splm-mis_splm-mis-minio-data
driver: local
splm-mis-postgres-data:
name: splm-mis_splm-mis-postgres-data
driver: local
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ ls -la
total 192
drwxr-xr-x 11 lahat lahat  4096 Aug  6 23:32 .
drwxr-xr-x  4 root  root   4096 Aug  6 21:40 ..
drwxrwxr-x  3 lahat lahat  4096 Aug  6 21:39 .agents
drwxrwxr-x  3 lahat lahat  4096 Aug  6 21:39 .claude
-rw-rw-r--  1 lahat lahat   737 Aug  6 21:39 .dockerignore
-rw-rw-r--  1 lahat lahat  2974 Aug  6 21:39 .env.example
-rw-------  1 root  root   1797 Aug  6 23:53 .env.prod
-rw-------  1 root  root   1512 Aug  6 21:50 .env.prod.bak.20260806-233202
-rw-rw-r--  1 lahat lahat  4912 Aug  6 21:39 .env.prod.example
drwxrwxr-x  7 lahat lahat  4096 Aug  7 10:55 .git
-rw-rw-r--  1 lahat lahat    38 Aug  6 21:39 .gitattributes
drwxrwxr-x  3 lahat lahat  4096 Aug  6 21:39 .github
-rw-rw-r--  1 lahat lahat  1259 Aug  6 21:39 .gitignore
drwxrwxr-x  3 lahat lahat  4096 Aug  6 21:39 .mvn
-rw-rw-r--  1 lahat lahat  3208 Aug  6 21:39 DEPLOY.md
-rw-rw-r--  1 lahat lahat  3395 Aug  6 21:39 Dockerfile
drwx------  2 1001  1001   4096 Aug  6 21:50 certs
-rw-rw-r--  1 lahat lahat  4064 Aug  6 21:39 docker-compose.prod.yml
-rw-rw-r--  1 lahat lahat  1939 Aug  6 21:39 docker-compose.yml
-rw-rw-r--  1 lahat lahat 11790 Aug  6 21:39 mvnw
-rw-rw-r--  1 lahat lahat  8481 Aug  6 21:39 mvnw.cmd
drwxrwxr-x  4 lahat lahat  4096 Aug  6 21:39 nginx
-rw-rw-r--  1 lahat lahat 62322 Aug  6 21:39 package-lock.json
-rw-rw-r--  1 lahat lahat   809 Aug  6 21:39 package.json
-rw-rw-r--  1 lahat lahat  7028 Aug  6 21:39 pom.xml
drwxrwxr-x  2 lahat lahat  4096 Aug  6 21:39 scripts
drwxrwxr-x  4 lahat lahat  4096 Aug  6 21:39 src
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$ git status
warning: could not open directory 'certs/': Permission denied
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
lahat@ubuntu-4gb-nbg1-1:/opt/splm-mis$