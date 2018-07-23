db.services.find({}).forEach(function (item){
    // Update hosts array
    item.hosts = ['localhost:8070']
    
    // Update default host and port
    item.host = 'localhost'
    item.port = '8070'
    
    // Update documentation location
    item.apidocumentation = "http://localhost:8070/api/experience/eresults/console"
    
    // Update healthcheck url - replace env=ci with the env
    item.healthcheckurl = "http://localhost:8001/amc/mule/applications/gplatform.experience.eresults-1.0.0/healthcheck?env=ci"
    
    // Update api management host and port
    item.servicemanagementhost = 'localhost'
    item.servicemanagementport = '8001'
    
    // Update all api management endpoints
    item.servicemanagementendpoints['redeploy'] = item.servicemanagementendpoints['redeploy'].replace('DEV', 'CI')
    item.servicemanagementendpoints['restart'] = item.servicemanagementendpoints['restart'].replace('DEV', 'CI')
    item.servicemanagementendpoints['undeploy'] = item.servicemanagementendpoints['undeploy'].replace('DEV', 'CI')
    item.servicemanagementendpoints['backup'] = item.servicemanagementendpoints['backup'].replace('DEV', 'CI')
    item.servicemanagementendpoints['logs'] = item.servicemanagementendpoints['logs'].replace('DEV', 'CI')
    item.servicemanagementendpoints['logs'] = item.servicemanagementendpoints['logs'].replace('logs', 'log')
    db.services.save(item)
})