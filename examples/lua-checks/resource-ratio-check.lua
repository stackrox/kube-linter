-- Example Lua check: Validate memory to CPU ratio
-- This demonstrates arithmetic operations and custom business logic

function check()
    local diagnostics = {}
    
    -- Extract pod specification
    local podSpec = extract.podSpec()
    if not podSpec then
        return diagnostics
    end
    
    if not podSpec.containers then
        return diagnostics
    end
    
    for i, container in ipairs(podSpec.containers) do
        if container.resources and container.resources.requests then
            local requests = container.resources.requests
            local cpu = requests.cpu
            local memory = requests.memory
            
            if cpu and memory then
                -- Simple heuristic: parse memory in bytes and CPU in millicores
                local memoryBytes = parseMemory(memory)
                local cpuMillis = parseCPU(cpu)
                
                if memoryBytes > 0 and cpuMillis > 0 then
                    -- Calculate memory per CPU core in GB
                    local memoryPerCore = (memoryBytes / (cpuMillis / 1000)) / (1024 * 1024 * 1024)
                    
                    -- Flag if memory to CPU ratio is too high (> 8GB per core)
                    if memoryPerCore > 8 then
                        table.insert(diagnostics, diagnostic(
                            string.format(
                                "Container '%s' has high memory/CPU ratio: %.2fGB per core. " ..
                                "Consider rebalancing resources.",
                                container.name or "unknown",
                                memoryPerCore
                            )
                        ))
                    end
                end
            end
        end
    end
    
    return diagnostics
end

-- Helper function to parse Kubernetes memory strings to bytes
function parseMemory(memStr)
    if not memStr then return 0 end
    
    local num, unit = string.match(memStr, "^(%d+%.?%d*)(%a*)$")
    if not num then return 0 end
    
    num = tonumber(num)
    if not num then return 0 end
    
    unit = unit or ""
    
    if unit == "Ki" then return num * 1024
    elseif unit == "Mi" then return num * 1024 * 1024
    elseif unit == "Gi" then return num * 1024 * 1024 * 1024
    elseif unit == "Ti" then return num * 1024 * 1024 * 1024 * 1024
    elseif unit == "k" then return num * 1000
    elseif unit == "M" then return num * 1000 * 1000
    elseif unit == "G" then return num * 1000 * 1000 * 1000
    elseif unit == "T" then return num * 1000 * 1000 * 1000 * 1000
    else return num end
end

-- Helper function to parse Kubernetes CPU strings to millicores
function parseCPU(cpuStr)
    if not cpuStr then return 0 end
    
    -- Handle millicores (e.g., "100m")
    local num = string.match(cpuStr, "^(%d+)m$")
    if num then return tonumber(num) end
    
    -- Handle cores (e.g., "1", "0.5")
    num = tonumber(cpuStr)
    if num then return num * 1000 end
    
    return 0
end