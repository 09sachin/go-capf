--Material views queries

CREATE MATERIALIZED VIEW workflow_table AS
WITH rankedworkflows AS (                                                                                                         
          SELECT tms_m_case_workflow_new.workflow_id,                                                                               
             tms_m_case_workflow_new.process_desc,                                                                                  
             row_number() OVER (PARTITION BY tms_m_case_workflow_new.workflow_id ORDER BY tms_m_case_workflow_new.crt_dt) AS row_num
            FROM tms_m_case_workflow_new                                                                                            
           WHERE ((tms_m_case_workflow_new.state_code)::text = '91'::text)                                                          
         )                                                                                                                          
  SELECT rankedworkflows.workflow_id,                                                                                               
     rankedworkflows.process_desc,                                                                                                  
     rankedworkflows.row_num                                                                                                        
    FROM rankedworkflows                                                                                                            
   WHERE (rankedworkflows.row_num = 1);



CREATE MATERIALIZED VIEW queries AS
  SELECT wa.remarks,                                                                                                                                                                                                        
     reim.card_no,                                                                                                                                                                                                          
     reim.claim_sub_dt,                                                                                                                                                                                                     
     reim.case_no,                                                                                                                                                                                                          
     wa.crt_dt                                                                                                                                                                                                              
    FROM (tms_t_case_workflow_audit wa                                                                                                                                                                                      
      JOIN case_dump_capf_reim_pfms reim ON (((reim.patient_no)::text = (wa.transaction_id)::text)))                                                                                                                        
   WHERE (((wa.current_group_id)::text = ANY ((ARRAY['GP603'::character varying, 'GPSHA'::character varying, 'GPMD'::character varying, 'GPBANK'::character varying])::text[])) AND ((reim.ben_pending)::text = 'Y'::text));




CREATE MATERIALIZED VIEW claims AS
 SELECT DISTINCT                                                        
     rem.case_no,                                                                              
     rem.claim_sub_dt,                                                                         
     rw.process_desc,                                                                          
     rem.claim_sub_amt,                                                                        
     rem.claim_app_amt,                                                                        
     rem.claim_paid_amt,                                                                       
     rem.workflow_id,                                                                          
     ttp.hosp_name,
     rem.card_no                                                                            
    FROM case_dump_capf_reim_pfms rem                                                       
      LEFT JOIN tms_t_reimbursement ttp ON rem.patient_no::text = ttp.patient_no::text
      JOIN workflow_table rw ON rem.workflow_id = rw.workflow_id;




CREATE MATERIALIZED VIEW track_case AS
  SELECT reimb.case_no,                                                                               
     reimb.claim_sub_dt,                                                                              
     workflow.process_desc,                                                                           
     wa.crt_dt, 
     wa.remarks, 
     wa.amount,                                                                                    
     reimb.card_no                                                                                    
    FROM ((tms_t_case_workflow_audit wa                                                               
      JOIN case_dump_capf_reim_pfms reimb ON (((wa.transaction_id)::text = (reimb.patient_no)::text)))
      JOIN workflow_table workflow ON ((wa.next_workflow_id = workflow.workflow_id)));


CREATE MATERIALIZED VIEW hospitals AS
  SELECT * from hem_t_hosp_info;