-- Example Lua check: Detect containers using both privileged mode and host networking
-- This demonstrates a complex validation that requires checking multiple conditions
-- across different parts of the pod specification

function check()
    local diagnostics = {}
    
    -- Extract pod specification
    local podSpec = extract.podSpec()
    if not podSpec then
        return diagnostics
    end
    
    -- Check if any container is privileged
    local hasPrivileged = false
    if podSpec.containers then
        for i, container in ipairs(podSpec.containers) do
            if container.securityContext and 
               container.securityContext.privileged == true then
                hasPrivileged = true
                break
            end
        end
    end
    
    -- Also check init containers
    if not hasPrivileged and podSpec.initContainers then
        for i, container in ipairs(podSpec.initContainers) do
            if container.securityContext and 
               container.securityContext.privileged == true then
                hasPrivileged = true
                break
            end
        end
    end
    
    -- Check if pod uses host networking
    local usesHostNetwork = podSpec.hostNetwork == true
    
    -- If both conditions are true, add diagnostic
    if hasPrivileged and usesHostNetwork then
        table.insert(diagnostics, diagnostic(
            "Pod uses both privileged containers and host networking, " ..
            "which creates a significant security risk"
        ))
    end
    
    return diagnostics
end