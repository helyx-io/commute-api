select
    stf.stop_id, 
    stf.stop_name, 
    stf.stop_desc, 
    stf.stop_lat, 
    stf.stop_lon, 
    stf.location_type, 
    stf.arrival_time, 
    stf.departure_time, 
    stf.stop_sequence, 
    stf.direction_id, 
    stf.route_short_name, 
    stf.route_type, 
    stf.route_color, 
    stf.route_text_color, 
    stf.trip_id 
from
    gtfs_%s.stop_times_full stf inner join 
    gtfs_%s.calendars c on stf.service_id=c.service_id 
where
    stf.stop_id=%d and 
    c.start_date <= '%s' and 
    c.end_date >= '%s' and %s=true