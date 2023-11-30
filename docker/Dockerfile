# Start from scratch (empty image)
FROM scratch

# Copy the binary built from the local machine
COPY aip_food_lookup /.

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["/aip_food_lookup"]
