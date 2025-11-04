#!/bin/bash

echo "Starting Simple TestObject Controller..."

while true; do
    # Get all TestObjects
    ALL_OBJECTS=$(kubectl get testobjects -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.metadata.namespace}{" "}{.status.state}{"\n"}{end}')
    
    if [ -n "$ALL_OBJECTS" ]; then
        echo "$ALL_OBJECTS" | while IFS=' ' read -r name namespace state; do
            # Skip if name is empty
            if [ -z "$name" ]; then
                continue
            fi
            
            # If no status, set to Pending
            if [ -z "$state" ] || [ "$state" = "<no value>" ]; then
                echo "Setting $name to Pending"
                kubectl patch testobject "$name" -n "$namespace" --type='merge' --subresource=status -p='{"status":{"state":"Pending"}}'
            
            # If status is Pending, sleep 5 seconds and randomly set to Succeeded or Failed
            elif [ "$state" = "Pending" ]; then
                echo "Processing $name..."
                sleep 5
                
                # Generate random number (0 or 1)
                RANDOM_RESULT=$((RANDOM % 2))
                
                if [ $RANDOM_RESULT -eq 0 ]; then
                    # Success case
                    kubectl patch testobject "$name" -n "$namespace" --type='merge' --subresource=status -p='{"status":{"state":"Succeeded","message":"Task completed successfully!"}}'
                    echo "✅ $name succeeded"
                else
                    # Failure case  
                    kubectl patch testobject "$name" -n "$namespace" --type='merge' --subresource=status -p='{"status":{"state":"Failed","message":"Task failed due to unexpected error"}}'
                    echo "❌ $name failed"
                fi
            fi
        done
    fi
    
    sleep 2
done